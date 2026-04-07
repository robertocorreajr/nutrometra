package google

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// CalendarEvent represents an event to sync with Google Calendar.
type CalendarEvent struct {
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	TimeZone    string
}

// CalendarProvider defines operations for external calendar integration.
type CalendarProvider interface {
	CreateEvent(ctx context.Context, accessToken string, event CalendarEvent) (externalEventID string, err error)
	UpdateEvent(ctx context.Context, accessToken string, externalEventID string, event CalendarEvent) error
	DeleteEvent(ctx context.Context, accessToken string, externalEventID string) error
}

// GoogleCalendarProvider implements CalendarProvider using the Google Calendar API.
type GoogleCalendarProvider struct{}

// NewGoogleCalendarProvider creates a new GoogleCalendarProvider.
func NewGoogleCalendarProvider() *GoogleCalendarProvider {
	return &GoogleCalendarProvider{}
}

func (p *GoogleCalendarProvider) newService(ctx context.Context, accessToken string) (*calendar.Service, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	svc, err := calendar.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("google calendar: create service: %w", err)
	}
	return svc, nil
}

func toGoogleEvent(event CalendarEvent) *calendar.Event {
	tz := event.TimeZone
	if tz == "" {
		tz = "America/Sao_Paulo"
	}
	return &calendar.Event{
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		Start: &calendar.EventDateTime{
			DateTime: event.StartTime.Format(time.RFC3339),
			TimeZone: tz,
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndTime.Format(time.RFC3339),
			TimeZone: tz,
		},
	}
}

// CreateEvent creates a new event on the user's primary Google Calendar.
func (p *GoogleCalendarProvider) CreateEvent(ctx context.Context, accessToken string, event CalendarEvent) (string, error) {
	svc, err := p.newService(ctx, accessToken)
	if err != nil {
		return "", err
	}
	gEvent := toGoogleEvent(event)
	created, err := svc.Events.Insert("primary", gEvent).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("google calendar: create event: %w", err)
	}
	return created.Id, nil
}

// UpdateEvent updates an existing event on the user's primary Google Calendar.
func (p *GoogleCalendarProvider) UpdateEvent(ctx context.Context, accessToken string, externalEventID string, event CalendarEvent) error {
	svc, err := p.newService(ctx, accessToken)
	if err != nil {
		return err
	}
	gEvent := toGoogleEvent(event)
	_, err = svc.Events.Update("primary", externalEventID, gEvent).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("google calendar: update event %s: %w", externalEventID, err)
	}
	return nil
}

// DeleteEvent deletes an event from the user's primary Google Calendar.
func (p *GoogleCalendarProvider) DeleteEvent(ctx context.Context, accessToken string, externalEventID string) error {
	svc, err := p.newService(ctx, accessToken)
	if err != nil {
		return err
	}
	err = svc.Events.Delete("primary", externalEventID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("google calendar: delete event %s: %w", externalEventID, err)
	}
	return nil
}
