package soundFX

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"io"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
)

// AudioConfig defines audio metadata that can be loaded from JSON
type AudioConfig struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Volume   float64 `json:"volume"`
	Category string  `json:"category"` // e.g., "sfx", "music", "voice"
}

// AudioDatabase holds all audio configurations
type AudioDatabase struct {
	Audio []AudioConfig `json:"audio"`
}

type AutomatedAudioLoader struct {
	loader    *resource.Loader
	audioMap  map[string]resource.AudioID
	idCounter resource.AudioID
	config    AudioDatabase
}

func NewAutomatedAudioLoader() *AutomatedAudioLoader {
	audioContext := audio.NewContext(44100)
	loader := resource.NewLoader(audioContext)

	loader.OpenAssetFunc = func(path string) io.ReadCloser {
		data, err := assets.SoundDir.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read asset %s: %v", path, err)
		}
		return io.NopCloser(bytes.NewReader(data))
	}

	return &AutomatedAudioLoader{
		loader:    loader,
		audioMap:  make(map[string]resource.AudioID),
		idCounter: 1, // Start from 1 (0 is typically reserved for "None")
	}
}

// Method 1: Load from JSON configuration file
func (aal *AutomatedAudioLoader) LoadFromConfig(configPath string) error {
	data, err := assets.DataDir.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	if err := json.Unmarshal(data, &aal.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Convert config to resource registry
	audioResources := make(map[resource.AudioID]resource.AudioInfo)

	for _, audioConfig := range aal.config.Audio {
		id := aal.idCounter
		aal.idCounter++

		audioResources[id] = resource.AudioInfo{
			Path:   audioConfig.Path,
			Volume: audioConfig.Volume,
		}

		aal.audioMap[audioConfig.Name] = id
	}

	aal.loader.AudioRegistry.Assign(audioResources)
	return nil
}

func (aal *AutomatedAudioLoader) PlayByName(name string) error {
	id, exists := aal.audioMap[name]
	if !exists {
		return fmt.Errorf("audio '%s' not found", name)
	}

	audio := aal.loader.LoadAudio(id)
	if err := audio.Player.Rewind(); err != nil {
		return fmt.Errorf("failed to rewind audio: %w", err)
	}

	audio.Player.Play()
	return nil
}

func (aal *AutomatedAudioLoader) GetAudioID(name string) (resource.AudioID, bool) {
	id, exists := aal.audioMap[name]
	return id, exists
}

func (aal *AutomatedAudioLoader) ListAudio() []string {
	names := make([]string, 0, len(aal.audioMap))
	for name := range aal.audioMap {
		names = append(names, name)
	}
	return names
}

// Helper functions
func isAudioFile(ext string) bool {
	audioExts := []string{".wav", ".ogg", ".mp3", ".m4a", ".flac"}
	for _, audioExt := range audioExts {
		if ext == audioExt {
			return true
		}
	}
	return false
}

func getDefaultVolume(path, name string) float64 {
	// Apply volume rules based on conventions
	lowerPath := strings.ToLower(path)
	lowerName := strings.ToLower(name)

	// Music should be quieter
	if strings.Contains(lowerPath, "music") || strings.Contains(lowerPath, "bgm") {
		return -0.2
	}

	// UI sounds should be moderate
	if strings.Contains(lowerPath, "ui") || strings.Contains(lowerName, "button") {
		return -0.1
	}

	// Explosion/impact sounds can be louder
	if strings.Contains(lowerName, "explosion") || strings.Contains(lowerName, "crash") {
		return 0.1
	}

	// Default volume
	return 0.0
}

/*// Usage example with loading screen
func ExampleUsage() {
	audioLoader := NewAutomatedAudioLoader()

	// Method 1: Load from config file
	if err := audioLoader.LoadFromConfig("assets/audio/config.json"); err != nil {
		// Fallback to auto-discovery
		log.Printf("Config loading failed, using auto-discovery: %v", err)
		if err := audioLoader.AutoDiscoverAudio("assets/audio"); err != nil {
			log.Fatalf("Failed to load audio: %v", err)
		}
	}

	// Show loading screen with progress
	audioLoader.PreloadAllWithProgress(func(current, total int, name string) {
		fmt.Printf("Loading audio: %s (%d/%d)\n", name, current, total)
		// Update your loading screen here
	})

	// Later in your game, play sounds by name
	if err := audioLoader.PlayByName("jump"); err != nil {
		log.Printf("Failed to play jump sound: %v", err)
	}

	// Or preload only specific categories
	audioLoader.PreloadByCategory("sfx") // Load only sound effects
}
*/
