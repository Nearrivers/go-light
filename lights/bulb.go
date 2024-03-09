package lights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/nearrivers/go-light/app_errors"
	"github.com/nearrivers/go-light/request"
)

type owner struct {
	Rid   string `json:"rid"`
	Rtype string `json:"rtype"`
}

type on struct {
	On bool `json:"on"`
}

type dimming struct {
	Brightness  float32 `json:"brightness"`
	MinDimLevel float32 `json:"min_dim_level"`
}

type colorGamut struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type gamut struct {
	Red   colorGamut `json:"red"`
	Green colorGamut `json:"green"`
	Blue  colorGamut `json:"blue"`
}

type color struct {
	Xy    colorGamut `json:"xy"`
	Gamut gamut      `json:"gamut"`
	// Types de gammes supportées par Hue
	//
	// A -> Gamme des produits Philips uniquement en couleur
	//
	// B -> Gamme limitée des premiers produits sortis en couleur
	//
	// C -> Gamme des produits Hue blanc et autres couleurs d'ambiances
	GamutType string `json:"gamut_type"`
}

type LightGet struct {
	Id       string   `json:"id"`
	Owner    owner    `json:"owner"`
	Metadata metaData `json:"metadata"`
	On       on       `json:"on"`
	Dimming  dimming  `json:"dimming"`
	Color    color    `json:"color"`
}

func GetOneBulbState(bulbId string) ([]LightGet, error) {
	resp, err := request.NewHueBodylessRequest(request.GET, "clip/v2/resource/light")
	if err != nil {
		return []LightGet{}, app_errors.RuntimeError{Err: err}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []LightGet{}, app_errors.RuntimeError{Err: err}
	}

	lightGet := struct {
		Errors []app_errors.HueApiError `json:"errors"`
		Data   []LightGet               `json:"data"`
	}{}

	err = json.Unmarshal(body, &lightGet)
	if err != nil {
		return []LightGet{}, app_errors.RuntimeError{Err: err}
	}

	return lightGet.Data, nil
}

func GetBulbsState() ([]LightGet, error) {
	resp, err := request.NewHueBodylessRequest(request.GET, "clip/v2/resource/light")
	if err != nil {
		return []LightGet{}, app_errors.RuntimeError{Err: err}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []LightGet{}, app_errors.RuntimeError{Err: err}
	}

	lightGet := struct {
		Errors []app_errors.HueApiError `json:"errors"`
		Data   []LightGet               `json:"data"`
	}{}

	err = json.Unmarshal(body, &lightGet)
	if err != nil {
		return []LightGet{}, app_errors.RuntimeError{Err: err}
	}

	if len(lightGet.Errors) > 0 {
		fmt.Println(lightGet.Errors)
	}

	return lightGet.Data, nil
}

func ControlLights() error {
	bulbs, err := GetBulbsState()
	if err != nil {
		return err
	}

	for _, bulb := range bulbs {
		newBulbState, err := json.Marshal(struct {
			Type string `json:"type"`
			On   on     `json:"on"`
		}{
			Type: "light",
			On: on{
				On: !bulb.On.On,
			},
		})

		if err != nil {
			return err
		}

		body := bytes.NewBuffer(newBulbState)
		resp, err := request.NewHueRequestWithBody(request.PUT, fmt.Sprintf("%s%s", "clip/v2/resource/light/", bulb.Id), body)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		responseBody := struct {
			Errors []app_errors.HueApiError `json:"errors"`
			Data   []struct {
				Rid   string `json:"rid"`
				Rtype string `json:"rtype"`
			} `json:"data"`
		}{}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &responseBody)
		if err != nil {
			return err
		}
	}

	return nil
}
