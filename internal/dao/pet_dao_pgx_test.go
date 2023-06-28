package dao_test

import (
	"log"
	"testing"

	"github.com/irwin13/go-petstore/pkg/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PetDaoSuite struct {
	suite.Suite
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *PetDaoSuite) SetupSuite() {
	log.Println("SetupSuite running")
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *PetDaoSuite) TearDownSuite() {
	log.Println("TearDownSuite running")
}

// The SetupTest method will be run before every test in the suite.
func (suite *PetDaoSuite) SetupTest() {
	log.Println("SetupTest running")
	truncateSql := []string{"TRUNCATE TABLE pet"}
	err := executeRawSql(truncateSql)
	assert := assert.New(suite.T())
	assert.Nil(err)
}

// The TearDownTest method will be run after every test in the suite.
func (suite *PetDaoSuite) TearDownTest() {
	log.Println("TearDownTest running")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPetDaoSuite(t *testing.T) {
	suite.Run(t, new(PetDaoSuite))
}

func (suite *PetDaoSuite) TestPetDao_Crud() {
	assert := assert.New(suite.T())

	// insert
	name := "name"
	description := "desc"

	insertModel := entity.InsertPet{}
	insertModel.Name = &name
	insertModel.Description = &description

	ID, err := petDao.Insert(insertModel)

	assert.Nil(err)
	assert.NotEmpty(ID)

	// search
	list, err := petDao.Search(ID)
	assert.Nil(err)
	assert.Len(list, 1, "search result not 1")

	// update
	updateName := "update"

	updateModel := entity.UpdatePet{}
	updateModel.ID = &ID
	updateModel.Name = &updateName
	updateModel.Description = &description

	result, err := petDao.Update(updateModel)
	assert.Nil(err)
	assert.Equal(int64(1), result, "update result not equal")

	list, err = petDao.Search(ID)
	assert.Nil(err)

	pet := list[0]
	assert.Equal(updateName, *pet.Name, "field name after update not equal")

	// delete
	deleteRequest := entity.DeletePet{
		ID: &ID,
	}
	result, err = petDao.Delete(deleteRequest)
	assert.Nil(err)
	assert.Equal(int64(1), result, "delete result not equal")

	list, err = petDao.Search(ID)
	assert.Nil(err)
	assert.Equal(0, len(list))

}
