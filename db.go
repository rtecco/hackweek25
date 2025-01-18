package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/jmoiron/sqlx"
)

type Profile struct {
	ID                  int    `db:"id"`
	Name                string `db:"name"`
	ProfileSVG          string `db:"profile_image_svg"`
	Summary             string `db:"summary"`
	SpecialtyTagsRaw    string `db:"specialty_tags"`
	SpecialtyTags       []string
	Vouches             int    `db:"vouches"`
	PortfolioDescriptor string `db:"portfolio_descriptor"`
}

func getProfile(db *sqlx.DB, id int) Profile {

	profile := Profile{}
	err := db.Get(&profile, "SELECT * FROM profiles WHERE id=?", id)

	if err != nil {
		log.Fatal(err)
	}

	profile.SpecialtyTags = strings.Split(profile.SpecialtyTagsRaw, ",")

	return profile
}

func savePortfolioDescriptor(index bleve.Index, db *sqlx.DB, id int, descriptor string) {

	_, err := db.Exec("UPDATE profiles SET portfolio_descriptor = ? WHERE id = ?", descriptor, id)

	if err != nil {
		log.Fatal(err)
	}

	buildIndex(index, db)
}

func saveSpecialtyTag(index bleve.Index, db *sqlx.DB, id int, tag string) {
	tag = strings.TrimSpace(tag)

	_, err := db.Exec("UPDATE profiles SET specialty_tags = specialty_tags || ',' || ? WHERE id = ?", tag, id)

	if err != nil {
		log.Fatal(err)
	}

	buildIndex(index, db)
}

func newIndex() bleve.Index {

	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)

	if err != nil {
		log.Fatal(err)
	}

	return index
}

func buildIndex(index bleve.Index, db *sqlx.DB) {

	profiles := []Profile{}
	db.Select(&profiles, "SELECT * FROM profiles")

	for _, profile := range profiles {
		fmt.Println("indexing", profile)
		index.Index(strconv.Itoa(profile.ID), &profile)
	}
}

func search(index bleve.Index, db *sqlx.DB, queryString string) []Profile {

	query := bleve.NewQueryStringQuery(queryString)
	search := bleve.NewSearchRequest(query)

	results, err := index.Search(search)

	if err != nil {
		log.Fatal(err)
	}

	profiles := []Profile{}

	for _, hit := range results.Hits {
		id, err := strconv.Atoi(hit.ID)

		if err != nil {
			log.Fatal(err)
		}

		if hit.Score < 0.005 {
			continue
		}

		profile := getProfile(db, id)
		profiles = append(profiles, profile)
	}

	return profiles
}

const schemaStmt = `
DROP TABLE IF EXISTS profiles;

CREATE TABLE IF NOT EXISTS profiles (id INTEGER NOT NULL PRIMARY KEY,
	name TEXT,
	profile_image_svg TEXT,
	summary TEXT,
	specialty_tags TEXT,
	vouches INT,
	portfolio_descriptor TEXT);

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Sarah Martinez', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#FFE0B2"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 75 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 25 35 Q 50 25 75 35" fill="#333" stroke="none"/>
  <rect x="30" y="40" width="40" height="5" fill="#555"/>
</svg>', 'Lawyer @ Pinnacle & Associates', 'corporate litigation', 2, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Michael Chen', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#F5D0C5"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 70 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 20 30 L 80 30 L 50 10 Z" fill="#333"/>
</svg>', 'Architect @ Skyline Innovations', 'sustainable design, urban high-rise, LEED certification, green remodels', 1, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Emma Thompson', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#FFDDC1"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 75 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 15 50 Q 50 80 85 50 Q 50 20 15 50" fill="#8B4513"/>
</svg>', 'Therapist @ Mindful Horizons', 'cognitive behavioral therapy, trauma recovery, EMDR', 0, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('James Wilson', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#E3D0C3"/>
  <rect x="25" y="40" width="50" height="10" rx="5" fill="#555"/>
  <circle cx="35" cy="45" r="5" fill="#FFF"/>
  <circle cx="65" cy="45" r="5" fill="#FFF"/>
  <path d="M 40 65 Q 50 70 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 25 25 Q 50 20 75 25 L 75 35 Q 50 30 25 35 Z" fill="#333"/>
</svg>', 'Coder @ ByteCraft Solutions', 'kubernetes, microservices, golang, distributed systems', 0, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Maria Garcia', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#F4C2A6"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 72 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 20 35 Q 50 60 80 35 Q 50 0 20 35" fill="#472203"/>
</svg>', 'Lawyer @ Justice Dynamics', 'intellectual property, patent law, trademark litigation', 1, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('David Kim', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#F5D5C5"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 68 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 25 30 Q 50 25 75 30 L 75 40 Q 50 35 25 40 Z" fill="#111"/>
</svg>', 'Coder @ Quantum Code Labs', 'machine learning, pytorch, computer vision', 1, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Jennifer Lee', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#FFE4C4"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 75 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 15 40 Q 50 70 85 40 Q 50 10 15 40" fill="#2C1810"/>
</svg>', 'Architect @ Urban Canvas Design', 'historic preservation, adaptive reuse, cultural facilities', 1, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Robert Taylor', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#E6C5B5"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 70 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 30 25 L 70 25 L 50 10 Z" fill="#555"/>
</svg>', 'Therapist @ Wellness Bridge', 'marriage counseling, attachment theory, gottman method, couples therapy', 0, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Lisa Wong', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#F5D0C5"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 72 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 20 40 Q 50 65 80 40 Q 50 15 20 40" fill="#111"/>
</svg>', 'Coder @ Digital Nexus Tech', 'react native, mobile development, app security', 1, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Daniel Park', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#F4D0C0"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 70 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 25 30 Q 50 25 75 30 L 75 35 Q 50 30 25 35 Z" fill="#333"/>
</svg>', 'Lawyer @ Legacy Law Group', 'estate planning, trust administration, elder law, probate', 0, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Rachel Anderson', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#FFE4C4"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 75 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 15 45 Q 50 75 85 45 Q 50 15 15 45" fill="#CD853F"/>
</svg>', 'Therapist @ Serenity Solutions', 'anxiety disorders, mindfulness-based therapy, eating disorders', 1, '');

INSERT INTO profiles (name, profile_image_svg, summary, specialty_tags, vouches, portfolio_descriptor)
VALUES ('Thomas Brown', '<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#E6C5B5"/>
  <circle cx="35" cy="45" r="5" fill="#555"/>
  <circle cx="65" cy="45" r="5" fill="#555"/>
  <path d="M 40 65 Q 50 70 60 65" stroke="#555" stroke-width="3" fill="none"/>
  <path d="M 20 30 L 80 30 L 50 15 Z" fill="#333"/>
</svg>', 'Architect @ Blueprint Masters', 'residential design, passive house', 2, '');
`
