package profiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Profile struct {
	Slug       string `json:"-"`
	RepoName   string `json:"repo"`
	BranchName string `json:"branch"`
	Dir        string `json:"dir"`
}

type Profiler struct {
	Profiles map[string]Profile
	file     *os.File
	decoder  *json.Decoder
	encoder  *json.Encoder
}

// Initializes a Profiler and loads profiler, any changes
// need to be saved using the Save() method.
//
// Do not forget to close the profiles file using the Close() method!
func New(profilePath string) (p *Profiler, err error) {
	f, err := os.OpenFile(profilePath, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			return
		}
		err = errors.Join(err, f.Close())
	}()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	Profiler := Profiler{
		file:     f,
		Profiles: make(map[string]Profile),
	}

	if stat.Size() == 0 {
		return &Profiler, nil
	}

	if err = Profiler.Load(); err != nil {
		return nil, fmt.Errorf("decoding profile: %w", err)
	}

	return &Profiler, nil
}

// Loads profiles from file into Profiles
func (p *Profiler) Load() error {
	clear(p.Profiles)

	if p.decoder == nil {
		p.decoder = json.NewDecoder(p.file)
	}
	if err := p.decoder.Decode(&p.Profiles); err != nil {
		return err
	}
	return nil
}

// Saves what is currently in Profiles to file
func (p *Profiler) Save() error {
	if _, err := p.file.Seek(0, 0); err != nil {
		return fmt.Errorf("seeking to start: %w", err)
	}
	if err := p.file.Truncate(0); err != nil {
		return fmt.Errorf("truncating file: %w", err)
	}

	if p.encoder == nil {
		p.encoder = json.NewEncoder(p.file)
		p.encoder.SetIndent("", "    ")
	}

	if err := p.encoder.Encode(p.Profiles); err != nil {
		return err
	}
	return nil
}

// Closes profiles file
func (p *Profiler) Close() error {
	return p.file.Close()
}

// Adds a profile, returns error if there already is a profile with slug.
func (p *Profiler) Add(slug string, profile Profile) error {
	if _, ok := p.Profiles[slug]; ok {
		return fmt.Errorf("profile '%s' already exists", slug)
	}
	p.Profiles[slug] = profile
	return nil
}

// Adds a profile, overwriting any existing one.
func (p *Profiler) AddOverwrite(slug string, profile Profile) {
	p.Profiles[slug] = profile
}

// Deletes profile and returns true if it deleted it.
func (p *Profiler) Delete(slug string) bool {
	_, ok := p.Profiles[slug]
	if ok {
		delete(p.Profiles, slug)
	}
	return ok
}

// Gets profile.
func (p *Profiler) Get(slug string) (Profile, bool) {
	profile, ok := p.Profiles[slug]
	profile.Slug = slug
	return profile, ok
}

// Cheks if profile exists
func (p *Profiler) Exists(slug string) bool {
	_, ok := p.Profiles[slug]
	return ok
}

func Get(slug string, profilePath string) (Profile, bool, error) {
	p, err := New(profilePath)
	if err != nil {
		return Profile{}, false, err
	}

	profile, ok := p.Get(slug)
	return profile, ok, p.Close()
}
