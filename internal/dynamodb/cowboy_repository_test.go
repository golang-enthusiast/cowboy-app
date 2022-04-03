package dynamodb

import (
	"testing"

	"cowboy-app/internal/domain"

	aws "github.com/aws/aws-sdk-go/aws"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gotest.tools/assert"
)

func TestCowboyFindByName(t *testing.T) {
	runSetup(t, loadTestCowboys)

	testCowboys := getTestCowboys()
	for _, cowboy := range testCowboys {
		found, err := testCowboyRepository.FindByName(cowboy.Name)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		assert.DeepEqual(t, *cowboy, *found)
	}
}

func TestCowboyList(t *testing.T) {
	runSetup(t, loadTestCowboys)

	testCowboys := getTestCowboys()

	// Find all.
	foundList, err := testCowboyRepository.List(&domain.CowboySearchCriteria{
		Limit: 100,
	})
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	assertCowboys(t, testCowboys, foundList)

	// Exclude "Peter".
	excludeName := "Peter"
	foundList, err = testCowboyRepository.List(&domain.CowboySearchCriteria{
		ExcludeName: excludeName,
		Limit:       100,
	})
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	excludePeterList := make([]*domain.Cowboy, 0)
	for _, cowboy := range testCowboys {
		if cowboy.Name != excludeName {
			excludePeterList = append(excludePeterList, cowboy)
		}
	}
	assertCowboys(t, excludePeterList, foundList)

	// Health gt 10.
	foundList, err = testCowboyRepository.List(&domain.CowboySearchCriteria{
		HealthGt: 10,
		Limit:    100,
	})
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	healthGt10List := make([]*domain.Cowboy, 0)
	for _, cowboy := range testCowboys {
		if cowboy.Health > 10 {
			healthGt10List = append(healthGt10List, cowboy)
		}
	}
	assertCowboys(t, healthGt10List, foundList)
}

func TestCowboyUpdateHealthPoints(t *testing.T) {
	runSetup(t, loadTestCowboys)

	testCowboys := getTestCowboys()
	for _, cowboy := range testCowboys {
		cowboy.Health = cowboy.Health - 1
		err := testCowboyRepository.UpdateHealthPoints(cowboy.Name, cowboy.Health)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		updatedCowboy, err := testCowboyRepository.FindByName(cowboy.Name)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
		assert.DeepEqual(t, *updatedCowboy, *cowboy)
	}
}

func loadTestCowboys() error {
	var (
		testCowboys = getTestCowboys()
		db          = awsDynamodb.New(testSession)
	)
	for _, cowboy := range testCowboys {
		attributes, err := dynamodbattribute.MarshalMap(cowboy)
		if err != nil {
			return err
		}
		input := &awsDynamodb.PutItemInput{
			Item:      attributes,
			TableName: aws.String(testCowboyTableName),
		}
		_, err = db.PutItem(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTestCowboys() []*domain.Cowboy {
	return []*domain.Cowboy{
		{
			Name:   "John",
			Health: 10,
			Damage: 1,
		},
		{
			Name:   "Bill",
			Health: 8,
			Damage: 2,
		},
		{
			Name:   "Sam",
			Health: 10,
			Damage: 1,
		},
		{
			Name:   "Peter",
			Health: 5,
			Damage: 3,
		},
		{
			Name:   "Philip",
			Health: 15,
			Damage: 1,
		},
	}
}

func assertCowboys(t *testing.T, expectedList, actualList []*domain.Cowboy) {
	var (
		expectedSize = len(expectedList)
		actualSize   = len(actualList)
	)
	if expectedSize != actualSize {
		t.Errorf("expected list size %d, but got %d", expectedSize, actualSize)
	}
	for _, expectedItem := range expectedList {
		for _, actualItem := range actualList {
			if expectedItem.Name == actualItem.Name {
				assert.DeepEqual(t, *expectedItem, *actualItem)
			}
		}
	}
}
