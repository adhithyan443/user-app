package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"user-app/config"
	"user-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/text/message"
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

	session := sessions.Default(ctx)

	msg := session.Get("message")
	session.Delete("message")
	session.Save()

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
		"user":    user,
		"message": msg,
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
	session := sessions.Default(ctx)

	if name == "" || email == "" || role == "" {
		session.Set("message", "All fields are required")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/edit/"+idParam)
		return
	}

	// Name validation
	if len(name) < 3 {
		session.Set("message", "Name must be at least 3 characters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/edit/"+idParam)
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRegex.MatchString(name) {
		session.Set("message", "Name should contain only letters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/edit/"+idParam)
		return
	}

	// Role validation
	if role != "admin" && role != "user" {
		session.Set("message", "Invalid role selected")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/edit/"+idParam)
		return
	}

	_, err = config.DB.Exec(
		"UPDATE users SET name=$1, email=$2, role=$3 WHERE id=$4",
		name, email, role, id,
	)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Update failed")
		return
	}

	session.Set("message", "User updated successfully")
	session.Save()

	ctx.Redirect(http.StatusFound, "/admin/users")
	// successMsg := "User updated successfully!"
	// ctx.Redirect(http.StatusFound, fmt.Sprintf("/admin/users/edit/%d?success=%s", id, successMsg))
}

func ShowUserPasswordPage(ctx *gin.Context) {

	session := sessions.Default(ctx)
	id := ctx.Param("id")
	msg := session.Get("message")
	session.Delete("message")
	session.Save()
	ctx.HTML(http.StatusOK, "admin_changepassword.html", gin.H{
		"message": msg,
		"id":      id,
	})
}

func EditUserPasswordPage(ctx *gin.Context) {

	session := sessions.Default(ctx)

	newpassword := ctx.PostForm("newpassword")
	confirmpassword := ctx.PostForm("confirmpassword")
	id := ctx.Param("id")

	if len(newpassword) < 6 {
		session.Set("message", "Password must be at least 6 characters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	if newpassword != confirmpassword {
		session.Set("message", "Password do not match")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)

	if err != nil {
		session.Set("message", "Failed to process password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	_, err = config.DB.Exec(
		"UPDATE users SET hashed_password=$1 WHERE id=$2",
		hashedPassword, id,
	)

	if err != nil {
		session.Set("message", "Failed to update password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	session.Set("message", "Password updated successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/admin/users")
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

	// Required fields
	if name == "" || email == "" || role == "" || password == "" {
		session.Set("message", "All fields are required")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	// Name validation (letters + space, min 3)
	if len(name) < 3 {
		session.Set("message", "Name must be at least 3 characters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRegex.MatchString(name) {
		session.Set("message", "Name should contain only letters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	// Password validation
	if len(password) < 6 {
		session.Set("message", "Password must be at least 6 characters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	// Role validation
	if role != "admin" && role != "user" {
		session.Set("message", "Invalid role selected")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

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
