package mojang

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

var (
	sessionAPI    = "https://sessionserver.mojang.com/session/minecraft/profile/"
	mojangNameAPI = "https://api.mojang.com/users/profiles/minecraft/"
)
var client *http.Client = http.DefaultClient

func SetMojangNameAPI(url string) {
	mojangNameAPI = url
}

func GetUUIDByNickname(nick string) (uuid.UUID, error) {
	resp, err := client.Get(mojangNameAPI + nick)
	if err != nil {
		return uuid.Nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return uuid.Nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var pr struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(pr.ID)
}

func GetHead64(id uuid.UUID, width, height int, overlay bool) (string, error) {
	sess, err := GetMojangProfile(id)
	if err != nil {
		return "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(sess.Properties[0].Value)
	if err != nil {
		return "", err
	}
	var tex texturesPayload
	if err := json.Unmarshal(decoded, &tex); err != nil {
		return "", err
	}

	skinB64, err := GetBase64FromURL(tex.Textures.SKIN.URL)
	if err != nil {
		return "", err
	}
	skinBytes, err := base64.StdEncoding.DecodeString(skinB64)
	if err != nil {
		return "", err
	}

	skinImg, _, err := image.Decode(bytes.NewReader(skinBytes))
	if err != nil {
		return "", err
	}

	head := image.NewRGBA(image.Rect(0, 0, 8, 8))
	draw.Draw(head, head.Bounds(), skinImg, image.Pt(8, 8), draw.Src)
	if overlay {
		layer := image.NewRGBA(image.Rect(0, 0, 8, 8))
		draw.Draw(layer, layer.Bounds(), skinImg, image.Pt(40, 8), draw.Over)
		draw.Draw(head, head.Bounds(), layer, image.Point{}, draw.Over)
	}

	resized := imaging.Resize(head, width, height, imaging.NearestNeighbor)

	var buf bytes.Buffer
	if err := png.Encode(&buf, resized); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func GetMojangProfile(id uuid.UUID) (*MojangSessionResponse, error) {
	plain := strings.ReplaceAll(id.String(), "-", "")
	resp, err := client.Get(sessionAPI + plain)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result MojangSessionResponse
	if body, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetBase64FromURL(url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}
