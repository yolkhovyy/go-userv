package oteld

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/baggage"
)

func ContextWithBaggageUserID(ctx context.Context, userID uuid.UUID) (context.Context, error) {
	memberUserID, err := baggage.NewMember("userID", userID.String())
	if err != nil {
		return ctx, fmt.Errorf("user id member new: %w", err)
	}

	bag, err := baggage.New(memberUserID)
	if err != nil {
		return ctx, fmt.Errorf("baggage new: %w", err)
	}

	ctx = baggage.ContextWithBaggage(ctx, bag)

	return ctx, nil
}

func ContextWithBaggageUserName(ctx context.Context, firstName, lastName string) (context.Context, error) {
	memberFirstName, err := baggage.NewMember("firstName", url.QueryEscape(firstName))
	if err != nil {
		return ctx, fmt.Errorf("first name member new: %w", err)
	}

	memberLastName, err := baggage.NewMember("lastName", url.QueryEscape(lastName))
	if err != nil {
		return ctx, fmt.Errorf("last name member new: %w", err)
	}

	bag, err := baggage.New(memberFirstName, memberLastName)
	if err != nil {
		return ctx, fmt.Errorf("baggage new: %w", err)
	}

	ctx = baggage.ContextWithBaggage(ctx, bag)

	return ctx, nil
}

func ContextWithBaggagePage(ctx context.Context, page, limit int, countryCode string) (context.Context, error) {
	memberPage, err := baggage.NewMember("page", strconv.Itoa(page))
	if err != nil {
		return ctx, fmt.Errorf("page member new: %w", err)
	}

	memberLimit, err := baggage.NewMember("limit", strconv.Itoa(limit))
	if err != nil {
		return ctx, fmt.Errorf("limit member new: %w", err)
	}

	memberCountryCode, err := baggage.NewMember("countryCode", countryCode)
	if err != nil {
		return ctx, fmt.Errorf("country code member new: %w", err)
	}

	bag, err := baggage.New(memberPage, memberLimit, memberCountryCode)
	if err != nil {
		return ctx, fmt.Errorf("baggage new: %w", err)
	}

	ctx = baggage.ContextWithBaggage(ctx, bag)

	return ctx, nil
}
