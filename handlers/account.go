package handlers

import (
	"errors"
	"regexp"
	"time"

	"github.com/bionicosmos/submgr/config"
	"github.com/bionicosmos/submgr/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var store *session.Store

func SessionInit() {
	store = session.New(session.Config{
		Expiration: time.Hour * 24 * 30,
		Storage: mongodb.New(mongodb.Config{
			ConnectionURI: config.Conf.DatabaseURL,
			Database:      config.Conf.DatabaseName,
		}),
	})
}

func SignInAccount(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	if session.Fresh() {
		account := new(models.Account)
		if err := c.BodyParser(account); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		inner, err := models.FindAccount(account.Username)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return fiber.NewError(fiber.StatusUnauthorized, "User does not exist.")
			}
			return err
		}
		if !account.PasswordIsCorrect(inner) {
			return fiber.ErrUnauthorized
		}
		store.CookieSessionOnly = !account.IsExtended
		session.Set("username", account.Username)
		session.Save()
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func SignUpAccount(c *fiber.Ctx) error {
	account := new(models.Account)
	if err := c.BodyParser(account); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	_, err := models.FindAccount(account.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err != mongo.ErrNoDocuments {
		return fiber.NewError(fiber.StatusConflict, "User exists.")
	}
	if err := account.Insert(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusCreated))
}

func SignOutAccount(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	if err := session.Destroy(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func ChangeAccountPassword(c *fiber.Ctx) error {
	password := new(struct {
		Old string `json:"old"`
		New string `json:"new"`
	})
	if err := c.BodyParser(password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	username, ok := session.Get("username").(string)
	if !ok {
		return errors.New("cannot get your username from the session")
	}
	inner, err := models.FindAccount(username)
	if err != nil {
		return err
	}
	account := &models.Account{Username: username, Password: password.Old}
	if !account.PasswordIsCorrect(inner) {
		return fiber.ErrUnauthorized
	}
	account.Password = password.New
	if err := account.Update(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func DeleteAccount(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	username, ok := session.Get("username").(string)
	if !ok {
		return errors.New("cannot get your username from the session")
	}
	if err := models.DeleteAccount(username); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func AuthorizeAccount(c *fiber.Ctx) error {
	subRegEx, err := regexp.Compile(`^\/api\/user\/[^/]+\/sub.*$`)
	if err != nil {
		return err
	}
	if path := c.Path(); path != "/api/account/sign-in" &&
		path != "/api/account/sign-up" &&
		!subRegEx.MatchString(path) {
		session, err := store.Get(c)
		if err != nil {
			return err
		}
		if session.Fresh() {
			return fiber.ErrUnauthorized
		}
	}
	return c.Next()
}
