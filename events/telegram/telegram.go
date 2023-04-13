package telegram

import (
	"bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/storage"
	"errors"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))
	for _, upd := range updates {
		res = append(res, event(upd))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Mesage:
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(update telegram.Update) events.Event {
	updType := fetchType(update)

	res := events.Event{
		Type: updType,
		Text: fetchText(update),
	}
	if updType == events.Mesage {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Mesage
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
