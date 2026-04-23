package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"user-app/config"
	"user-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUser(ctx *gin.Context) {

	session := sessions.Default(ctx)
	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	search := ctx.Query("search")

	var row *sql.Rows
	var err error
	if search != "" {
		query := `
			SELECT id,name,email,role 
			FROM users
			WHERE name ILIKE $1 OR email ILIKE $1
		`
		row, err = config.DB.Query(query, "%"+search+"%")
	} else {
		row, err = config.DB.Query("SELECT id,name,email,role FROM users")
	}

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error feching users")
		return
	}

	defer row.Close()

	var users []models.User

	for row.Next() {
		var user models.User

		err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			fmt.Println("Scan error:", err)
			continue
		}
		users = append(users, user)
	}

	ctx.HTML(http.StatusOK, "admin_users.html", gin.H{
		"users":   users,
		"message": msg,
	})
}

func EditUserPage(ctx *gin.Context) {

	idparam := ctx.Param("id")
	// fmt.Println("ID from URL:", idparam)
	id, err := strconv.Atoi(idparam)

	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid user ID")
		return
	}
	var user models.User

	err = config.DB.QueryRow(
		"SELECT id,name,email,role FROM users WHERE id=$1", id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "User not found")
		return
	}

	ctx.HTML(http.StatusOK, "edit_user.html", gin.H{
		"user": user,
	})
}

func UpdateUserPage(ctx *gin.Context) {

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	role := ctx.PostForm("role")

	_, err = config.DB.Exec(
		"UPDATE users SET name=$1, email=$2, role=$3 WHERE id=$4",
		name, email, role, id,
	)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Update failed")
		return
	}

	session := sessions.Default(ctx)
	session.Set("message", "User updated successfully")
	session.Save()

	ctx.Redirect(http.StatusFound, "/admin/users")
	// successMsg := "User updated successfully!"
	// ctx.Redirect(http.StatusFound, fmt.Sprintf("/admin/users/edit/%d?success=%s", id, successMsg))
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := config.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Delete failed")
		return
	}

	session := sessions.Default(ctx)
	session.Set("message", "User deleted successfully")
	session.Save()

	ctx.Redirect(http.StatusFound, "/admin/users")
}

func NewUserPage(ctx *gin.Context) {
	session := sessions.Default(ctx)

	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	ctx.HTML(http.StatusOK, "admin_add_user.html", gin.H{
		"message": msg,
	})
}

func AddNewUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	role := ctx.PostForm("role")
	password := ctx.PostForm("password")
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// ctx.HTML(http.StatusInternalServerError, "signup.html", gin.H{
		// 	"error": "Failed to hash password",
		// })
		session.Set("message", "Failed to hash password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	query := "INSERT INTO users(name,email,role, hashed_password) VALUES($1,$2,$3,$4)"
	_, err = config.DB.Exec(query, name, email, role, hashedpassword)

	if err != nil {
		// ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
		// 	"error": "Email already exists or invalid data",
		// })
		session.Set("message", "Email already exists or invalid data")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	session.Set("message", "User created successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/admin/users")
}
