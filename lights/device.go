package lights

import (
	"encoding/json"
	"io"

	"github.com/nearrivers/go-light/app_errors"
	"github.com/nearrivers/go-light/request"
)

// Désigne un appareil
type AppBulb struct {
	Rid        string  // Identifiant de l'appareil
	IsOn       bool    // Désigne l'état de l'ampoule. true pour allumé et false pour éteint
	Brightness float32 // Luminosité de l'ampoule
}

// Liste des appareils trouvés pendant l'exécution de l'application
type AppBulbs []AppBulb

type resourceIdentifierGet struct {
	Rid   string `json:"rid"`   // Identifiant unique de la ressource
	Rtype string `json:"rtype"` // Type de la ressource
}

type productData struct {
	ModelId              string `json:"model_id"`               // Identifiant unique du modèle de l'appareil
	ManufacturerName     string `json:"manufacturer_name"`      // Nom du fabricant de l'appareil
	ProductName          string `json:"product_name"`           // Nom du produit
	ProductArchetype     string `json:"product_archetype"`      // Archétype du produit
	Certified            bool   `json:"certified"`              // Cet appareil est reconnu par Hue
	SoftwareVersion      string `json:"software_version"`       // Version du logiciel de l'appareil
	HardwarePlatformType string `json:"hardware_platform_type"` // Type du matériel; identifié par le code du fabricant et ImageType
}

type metaData struct {
	Name      string `json:"name"`      // Nom de la ressource lisible par un humain
	Archetype string `json:"archetype"` // Par défaut égal à l'archétype donné par le fabricant. Modifiable par un utilisateur
}

type userTest struct {
	Status string `json:"status"` // Egal à "set" ou "changing"
	// Active ou prolonge le mode _usertest_ de l'utilisateur pour 120 secondes. false désactive le mode usertest.
	//
	// En mode usertest, les appareils signalent les changements d'états plus fréquemment et indiquent ces changements via, si elle existe, une LED présente physiquement sur l'appareil
	UserTest bool `json:"usertest"`
}

type DeviceGet struct {
	Type        string                  `json:"type"` // (device) Type des ressources supportées
	Id          string                  `json:"id"`   // Identifiant unique représentant une instance d'une ressource spécifique
	ProductData productData             `json:"product_data"`
	MetaData    metaData                `json:"metadata"` // Métadonnées
	UserTest    userTest                `json:"usertest"`
	Services    []resourceIdentifierGet `json:"services"` // Liste tous les services qui fournissent un c
}

// Récupère la liste des appareils Phillips Hue
func GetDeviceList() (*AppBulbs, []error) {
	resp, err := request.NewHueBodylessRequest(request.GET, "clip/v2/resource/device")

	if err != nil {
		return &AppBulbs{}, []error{app_errors.RuntimeError{Err: err}}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &AppBulbs{}, []error{app_errors.RuntimeError{Err: err}}
	}

	deviceGet := struct {
		Errors []app_errors.HueApiError `json:"errors"`
		Data   []DeviceGet              `json:"data"`
	}{}

	err = json.Unmarshal(body, &deviceGet)
	if err != nil {
		return &AppBulbs{}, nil
	}

	if len(deviceGet.Errors) > 0 {
		apiErrors := []error{}
		for _, err := range deviceGet.Errors {
			apiErrors = append(apiErrors, err)
		}

		return &AppBulbs{}, apiErrors
	}

	appBulbs := AppBulbs{}
	for _, device := range deviceGet.Data {
		for _, service := range device.Services {
			appBulbs = append(appBulbs, AppBulb{Rid: service.Rid})
		}
	}

	return &appBulbs, nil
}
