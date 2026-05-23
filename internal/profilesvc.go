package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
)

type ProfileServiceImpl struct {
	ProfilePath string
	httpClient  *http.Client
}

func NewProfileServiceImpl(profilePath string, httpClient *http.Client) *ProfileServiceImpl {
	return &ProfileServiceImpl{
		ProfilePath: profilePath,
		httpClient:  httpClient,
	}
}

func (p *ProfileServiceImpl) GetProfile(ctx context.Context, id int) (domain.Profile, error) {
	resp, err := p.httpClient.Get(p.ProfilePath + "/" + strconv.Itoa(id))
	if err != nil {
		return domain.Profile{}, err
	}
	defer resp.Body.Close()

	var profile domain.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return domain.Profile{}, err
	}

	return profile, nil
}
