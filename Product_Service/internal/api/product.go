package api

import (
	"Product_Service/internal/jwt"
	"Product_Service/internal/models"
	//"encoding/json"
	//"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (api *api) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "supplier" {
		http.Error(w, "Only suppliers can create products", http.StatusForbidden)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product.UserID = user.ID

	productID, err := api.db.CreateProduct(product)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusBadRequest)
		return
	}

	product.ID = productID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (api *api) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "supplier" {
		http.Error(w, "Only suppliers can delete products", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	productID, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(productID)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}
	err = api.db.DeleteProductByID(id)
	if err != nil {
		http.Error(w, "Error deleting a product", http.StatusInternalServerError)
		return
	}
}

func (api *api) GetAllProductsForClientHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "client" {
		http.Error(w, "Only clients can get all products", http.StatusForbidden)
		return
	}

	products, err := api.db.GetAllProductsForClient()
	if err != nil {
		http.Error(w, "Error getting products", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (api *api) GetAllProductsForSupplierHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "supplier" {
		http.Error(w, "Only suppliers can get all products", http.StatusForbidden)
		return
	}

	products, err := api.db.GetAllProductsForSupplier(user.ID)
	if err != nil {
		http.Error(w, "Error getting products", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (api *api) validateUserToken(r *http.Request) (*models.User, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}

	return user, nil
}
