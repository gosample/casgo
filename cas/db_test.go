package cas_test

import (
	. "github.com/t3hmrman/casgo/cas"

	//"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cas DB adapter", func() {

	Describe("DbExists function", func() {
		It("should return whether the database exists or not", func() {
			exists, casErr := testCASServer.Db.DbExists()
			Expect(casErr).To(BeNil())
			Expect(exists).To(Equal(true))
		})
	})

	Describe("GetDbName function", func() {
		It("should return the name of the server", func() {
			actual, expected := testCASServer.Db.GetDbName(), testCASServer.Config["dbName"]
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("GetUsersTableName function", func() {
		It("should return the default if not set differently", func() {
			actual, expected := testCASServer.Db.GetUsersTableName(), "users"
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("GetServicesTableName function", func() {
		It("should return the default if not set differently", func() {
			actual, expected := testCASServer.Db.GetServicesTableName(), "services"
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("GetTicketsTableName function", func() {
		It("should return the default if not set differently", func() {
			actual, expected := testCASServer.Db.GetTicketsTableName(), "tickets"
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("LoadJSONFixture function", func() {
		It("should not error when loading JSON into the database", func() {
			err := testCASServer.Db.LoadJSONFixture(
				testCASServer.Db.GetDbName(),
				testCASServer.Db.GetServicesTableName(),
				"fixtures/services.json",
			)
			Expect(err).To(BeNil())
		})

		It("Should increase the number of items in the given table", func() {
			err := testCASServer.Db.TeardownTable("services")
			Expect(err).To(BeNil())

			// Attempt to find a service in the fixture should fail
			service, err := testCASServer.Db.FindServiceByUrl("localhost:9090/validateCASLogin")
			Expect(err).ToNot(BeNil())
			Expect(service).To(BeNil())

			// Load the fixture
			err = testCASServer.Db.LoadJSONFixture(
				testCASServer.Db.GetDbName(),
				testCASServer.Db.GetServicesTableName(),
				"fixtures/services.json")
			Expect(err).To(BeNil())

			// Attempting to find a serive in the fixture shoudl pass now
			service, err = testCASServer.Db.FindServiceByUrl("localhost:9090/validateCASLogin")
			Expect(err).To(BeNil())
			Expect(service).ToNot(BeNil())

		})
	})

	Describe("FindServiceByUrl function", func() {
		It("should find the service added by the loaded test fixture", func() {
			// Create the service we expect to find in the fixture
			expectedService := &CASService{
				Name:       "test_service",
				Url:        "localhost:9090/validateCASLogin",
				AdminEmail: "noone@nowhere.com",
			}

			// Attempt to get a service by name
			returnedService, casErr := testCASServer.Db.FindServiceByUrl(expectedService.Url)
			Expect(casErr).To(BeNil())
			Expect(returnedService).ToNot(BeNil())
			Expect(returnedService).To(Equal(expectedService))
		})
	})

	Describe("FindUserByEmail function", func() {
		It("should find the user added by the loaded test fixture", func() {
			// Create the user we're expecting to get back
			expectedUser := &User{
				Email:    "test@test.com",
				Password: "thisisnotarealpassword",
			}

			// Attempt to get a user by name
			returnedUser, casErr := testCASServer.Db.FindUserByEmail(expectedUser.Email)
			Expect(casErr).To(BeNil())
			Expect(returnedUser).ToNot(BeNil())
			Expect(returnedUser).To(Equal(expectedUser))
		})
	})

	Describe("AddNewUser function", func() {
		It("should successfully add a new user", func() {
			// Add the user
			newUser, casErr := testCASServer.Db.AddNewUser("test_user@test.com", "randompassword")
			Expect(casErr).To(BeNil())
			Expect(newUser).ToNot(BeNil())

			// Find the user by email
			returnedUser, casErr := testCASServer.Db.FindUserByEmail(newUser.Email)
			Expect(casErr).To(BeNil())
			Expect(returnedUser).ToNot(BeNil())
			Expect(returnedUser.Email).To(Equal(newUser.Email))
		})
	})

	Describe("AddTicketForService function", func() {
		It("should successfully add a ticket for a given service", func() {
			// Create a new CASTicket to store
			ticket := &CASTicket{
				UserEmail:      "test@test.com",
				UserAttributes: map[string]string{},
				WasSSO:         false,
			}

			mockService := &CASService{
				Url:        "localhost:8080",
				Name:       "mock_service",
				AdminEmail: "noone@nowhere.com",
			}

			ticket, casErr := testCASServer.Db.AddTicketForService(ticket, mockService)
			Expect(casErr).To(BeNil())
			Expect(ticket).ToNot(BeNil())
			Expect(ticket.Id).ToNot(BeEmpty())
		})
	})

})
