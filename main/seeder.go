package main

import (
	"context"
	"github.com/google/martian/log"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"google.golang.org/api/iterator"
	"math/rand"
	"padel-backend/main/firebase"
	"padel-backend/main/match"
	"padel-backend/main/player"
	"padel-backend/main/util"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func Seed() {
	go PlayerSeeder(100)
	go MatchSeeder(5)

	wg.Wait()

	log.Infof("%v seeding complete!", util.GetLogPrefix("Seeder", "Seed"))
}

func PlayerSeeder(limit int) {
	wg.Add(1)

	var (
		ps []player.Player
		f  = faker.New()
	)

	for i := 0; i < limit; i++ {
		var (
			ID    = uuid.New().String()
			Name  = f.Person().Name()
			Email = f.Internet().Email()
		)

		ps = append(ps, player.Player{
			ID:    ID,
			Name:  Name,
			Email: Email,
		})
	}

	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	batch := firestore.Batch()

	for _, p := range ps {
		docRef := firestore.Collection(player.CollectionPlayer).NewDoc()

		_ = batch.Set(docRef, p)
	}

	_, err := batch.Commit(context.Background())
	if err != nil {
		log.Errorf("%v %v", util.GetLogPrefix("Seeder", "PlayerSeeder"), err)
	}

	wg.Done()
}

func MatchSeeder(limit int) {
	wg.Wait()
	wg.Add(1)

	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	var ps []player.Player

	iter := firestore.Collection(player.CollectionPlayer).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("%v %v", util.GetLogPrefix("Seeder", "MatchSeeder"), err)
		}

		var p player.Player
		err = doc.DataTo(&p)
		if err != nil {
			log.Errorf("%v %v", util.GetLogPrefix("Seeder", "MatchSeeder"), err)
		}

		ps = append(ps, p)
	}

	var (
		ms []match.Match
		f  = faker.New()
	)

	for i := 0; i < limit; i++ {
		var (
			ID       = uuid.New().String()
			Name     = f.Company().Name()
			Location = f.Address().StreetAddress()
			Time     = f.Time().Time(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC))
		)

		ms = append(ms, match.Match{
			ID:       ID,
			Club:     Name,
			Location: Location,
			Time:     Time,
			TeamA: match.Team{
				Player1: ps[rand.Intn(len(ps))],
				Player2: ps[rand.Intn(len(ps))],
				Results: []match.TeamResult{},
			},
			TeamB: match.Team{
				Player1: ps[rand.Intn(len(ps))],
				Player2: ps[rand.Intn(len(ps))],
				Results: []match.TeamResult{},
			},
			Locked: false,
		})
	}

	batch := firestore.Batch()

	for _, m := range ms {
		docRef := firestore.Collection(match.CollectionMatch).NewDoc()

		_ = batch.Set(docRef, m)
	}

	_, err := batch.Commit(context.Background())
	if err != nil {
		log.Errorf("%v %v", util.GetLogPrefix("Seeder", "MatchSeeder"), err)
	}

	wg.Done()
}
