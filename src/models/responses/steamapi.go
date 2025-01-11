package responses

type SteamResponse struct {
	Response struct {
		PlayerCount int `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

type SearchResult struct {
	Items []struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Price struct {
			Currency           string  `json:"currency"`
			Initial            float32 `json:"initial"`
			Final              float32 `json:"final"`
			DiscountPercentage int
		} `json:"price"`
		TinyImage    string `json:"tiny_image"`
		HeaderImage  string
		BigImage     string
		LibraryImage string
		Metascore    string `json:"metascore"`
		Platforms    struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		StreamingVideo    bool   `json:"streamingvideo"`
		ControllerSupport string `json:"controller_support"`
	} `json:"items"`
}

type SteamGameResponse struct {
	Success bool          `json:"success"`
	Data    SteamGameData `json:"data"`
}

// SteamGameData - структура, що містить дані про гру
type SteamGameData struct {
	Type                string         `json:"type"`
	Name                string         `json:"name"`
	SteamAppID          int            `json:"steam_appid"`
	IsFree              bool           `json:"is_free"`
	Price               *PriceOverview `json:"price_overview"`
	DLC                 []int          `json:"dlc"`
	DetailedDescription string         `json:"detailed_description"`
	AboutTheGame        string         `json:"about_the_game"`
	ShortDescription    string         `json:"short_description"`
	SupportedLanguages  string         `json:"supported_languages"`
	HeaderImage         string         `json:"header_image"`
	CapsuleImage        string         `json:"capsule_image"`
	CapsuleImageV5      string         `json:"capsule_imagev5"`
	CardImage           string
	Website             string             `json:"website"`
	Developers          []string           `json:"developers"`
	Publishers          []string           `json:"publishers"`
	Packages            []int              `json:"packages"`
	PackageGroups       []PackageGroup     `json:"package_groups"`
	Platforms           Platforms          `json:"platforms"`
	Categories          []Category         `json:"categories"`
	Genres              []Genre            `json:"genres"`
	Screenshots         []Screenshot       `json:"screenshots"`
	Movies              []Movie            `json:"movies"`
	Recommendations     Recommendations    `json:"recommendations"`
	Achievements        Achievements       `json:"achievements"`
	ReleaseDate         ReleaseDate        `json:"release_date"`
	SupportInfo         SupportInfo        `json:"support_info"`
	Background          string             `json:"background"`
	BackgroundRaw       string             `json:"background_raw"`
	ContentDescriptors  ContentDescriptors `json:"content_descriptors"`
	Ratings             Ratings            `json:"ratings"`
}

// PlatformRequirements - вимоги до платформи (Windows, Mac, Linux)
type PlatformRequirements struct {
	Minimum string `json:"minimum"`
}

// PackageGroup - структура для групи пакетів
type PackageGroup struct {
	Name                    string       `json:"name"`
	Title                   string       `json:"title"`
	Description             string       `json:"description"`
	SelectionText           string       `json:"selection_text"`
	SaveText                string       `json:"save_text"`
	IsRecurringSubscription string       `json:"is_recurring_subscription"`
	Subs                    []PackageSub `json:"subs"`
}

type PriceOverview struct {
	Currency           string `json:"currency"`
	Initial            int    `json:"initial"`
	Final              int    `json:"final"`
	DiscountPercentage int    `json:"discount_percent"`
	InitialFormatted   string `json:"initial_formatted"`
	FinalFormatted     string `json:"final_formatted"`
}

// PackageSub - структура для підписки в групі пакетів
type PackageSub struct {
	PackageID                int    `json:"packageid"`
	PercentSavingsText       string `json:"percent_savings_text"`
	PercentSavings           int    `json:"percent_savings"`
	OptionText               string `json:"option_text"`
	OptionDescription        string `json:"option_description"`
	CanGetFreeLicense        string `json:"can_get_free_license"`
	IsFreeLicense            bool   `json:"is_free_license"`
	PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
}

// Platforms - підтримувані платформи
type Platforms struct {
	Windows bool `json:"windows"`
	Mac     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

// Category - категорія гри
type Category struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// Genre - жанр гри
type Genre struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Screenshot - знімок екрану гри
type Screenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

// Movie - інформація про відео
type Movie struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	WebM      WebM   `json:"webm"`
	MP4       MP4    `json:"mp4"`
	Highlight bool   `json:"highlight"`
}

// WebM - структура для відео у форматі WebM
type WebM struct {
	_480 string `json:"480"`
	Max  string `json:"max"`
}

// MP4 - структура для відео у форматі MP4
type MP4 struct {
	_480 string `json:"480"`
	Max  string `json:"max"`
}

// Recommendations - рекомендації
type Recommendations struct {
	Total int `json:"total"`
}

// Achievements - досягнення
type Achievements struct {
	Total       int           `json:"total"`
	Highlighted []Achievement `json:"highlighted"`
}

// Achievement - окреме досягнення
type Achievement struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ReleaseDate - дата випуску гри
type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date"`
}

// SupportInfo - інформація про підтримку
type SupportInfo struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}

// ContentDescriptors - контентні дескриптори
type ContentDescriptors struct {
	IDs   []int  `json:"ids"`
	Notes string `json:"notes"`
}

// Ratings - рейтинги гри
type Ratings struct {
	USK          RatingDetail `json:"usk"`
	Agcom        RatingDetail `json:"agcom"`
	Cadpa        RatingDetail `json:"cadpa"`
	Dejus        RatingDetail `json:"dejus"`
	SteamGermany RatingDetail `json:"steam_germany"`
}

// RatingDetail - деталі рейтингу
type RatingDetail struct {
	Rating          string `json:"rating"`
	Descriptors     string `json:"descriptors"`
	RequiredAge     string `json:"required_age,omitempty"`
	Banned          string `json:"banned,omitempty"`
	UseAgeGate      string `json:"use_age_gate,omitempty"`
	RatingGenerated string `json:"rating_generated,omitempty"`
}
