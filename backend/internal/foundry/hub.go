package foundry

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Message struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type Conn struct {
	CampaignID uuid.UUID
	World      string
	Version    string
	ws         *websocket.Conn
	send       chan []byte
	pending    sync.Map
	closed     chan struct{}
	once       sync.Once
}

type Hub struct {
	mu    sync.RWMutex
	conns map[uuid.UUID]*Conn
}

func NewHub() *Hub {
	return &Hub{conns: map[uuid.UUID]*Conn{}}
}

func (h *Hub) Get(campaignID uuid.UUID) *Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.conns[campaignID]
}

func (h *Hub) set(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if old := h.conns[c.CampaignID]; old != nil {
		old.close()
	}
	h.conns[c.CampaignID] = c
}

func (h *Hub) remove(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns[c.CampaignID] == c {
		delete(h.conns, c.CampaignID)
	}
}

func (h *Hub) Serve(ws *websocket.Conn, campaignID uuid.UUID, onClose func()) {
	c := &Conn{
		CampaignID: campaignID,
		ws:         ws,
		send:       make(chan []byte, 16),
		closed:     make(chan struct{}),
	}
	h.set(c)
	go c.writer()
	c.readLoop()
	h.remove(c)
	c.close()
	if onClose != nil {
		onClose()
	}
}

func (c *Conn) writer() {
	for {
		select {
		case <-c.closed:
			return
		case b := <-c.send:
			if err := c.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				c.close()
				return
			}
		}
	}
}

func (c *Conn) readLoop() {
	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		var m Message
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		if m.ID != "" {
			if ch, ok := c.pending.Load(m.ID); ok {
				c.pending.Delete(m.ID)
				ch.(chan Message) <- m
				continue
			}
		}
		if m.Type == "hello" {
			var info struct {
				World   string `json:"world"`
				Version string `json:"version"`
			}
			_ = json.Unmarshal(m.Payload, &info)
			c.World = info.World
			c.Version = info.Version
		}
	}
}

func (c *Conn) close() {
	c.once.Do(func() {
		close(c.closed)
		_ = c.ws.Close()
	})
}

func (c *Conn) Request(ctx context.Context, typ string, payload any) (Message, error) {
	id := uuid.NewString()
	raw, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}
	msg, _ := json.Marshal(Message{ID: id, Type: typ, Payload: raw})

	ch := make(chan Message, 1)
	c.pending.Store(id, ch)
	defer c.pending.Delete(id)

	select {
	case c.send <- msg:
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case <-c.closed:
		return Message{}, errors.New("foundry connection closed")
	}

	select {
	case resp := <-ch:
		if resp.Error != "" {
			return resp, errors.New(resp.Error)
		}
		return resp, nil
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case <-c.closed:
		return Message{}, errors.New("foundry connection closed")
	}
}
