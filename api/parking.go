package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"br.com.mlabs/models"
	"br.com.mlabs/usecases"
	"br.com.mlabs/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// NewParkingRouter creates a subrouter for parking endpoints
func NewParkingRouter(router *mux.Router) {
	parkingRouter := router.PathPrefix("/parking").Subrouter()
	parkingRouter.HandleFunc("/{plate}", HistoryHandler).Methods("GET")
	parkingRouter.HandleFunc("/in", ImageRecognitionHandler).Methods("POST")
	parkingRouter.HandleFunc("/{id}/out", CheckoutHandler).Methods("PUT")
	parkingRouter.HandleFunc("/{id}/pay", PayHandler).Methods("PUT")
	parkingRouter.HandleFunc("", ReservationHandler).Methods("POST")
}

// ReservationHandler reserver a parking spot
func ReservationHandler(w http.ResponseWriter, r *http.Request) {
	var request models.ParkingRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(stringToJSON(utils.ErrBadRequest.Error()))

		return
	}

	id, err := usecases.MakeReservation(request)
	if err != nil {
		switch err {
		case utils.ErrInternalServer:
			w.WriteHeader(http.StatusInternalServerError)
		case utils.ErrPlateNotValid:
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(stringToJSON(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(idToJSON(id))
}

// HistoryHandler gets the plate's history
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	request := models.ParkingRequest{}

	vars := mux.Vars(r)
	request.Plate = vars["plate"]

	history, err := usecases.GetReservations(request)
	if err != nil {
		switch err {
		case utils.ErrPlateNotValid:
			w.WriteHeader(http.StatusBadRequest)
		case utils.ErrInternalServer:
			w.WriteHeader(http.StatusInternalServerError)
		case utils.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(stringToJSON(err.Error()))

		return
	}

	json, err := json.Marshal(history)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(stringToJSON(utils.ErrInternalServer.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// PayHandler sets the payment
func PayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := usecases.Pay(vars["id"]); err != nil {
		switch err {
		case utils.ErrIDNotValid:
			w.WriteHeader(http.StatusBadRequest)
		case utils.ErrInternalServer:
			w.WriteHeader(http.StatusInternalServerError)
		case utils.ErrAlreadyPaid:
			w.WriteHeader(http.StatusOK)
		case utils.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(stringToJSON(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(stringToJSON("Paid"))
}

// CheckoutHandler checks out a parking space
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := usecases.Checkout(vars["id"]); err != nil {
		switch err {
		case utils.ErrPayFirst:
			w.WriteHeader(http.StatusPaymentRequired)
		case utils.ErrIDNotValid:
			w.WriteHeader(http.StatusBadRequest)
		case utils.ErrAlreadyCheckedOut:
			w.WriteHeader(http.StatusOK)
		case utils.ErrInternalServer:
			w.WriteHeader(http.StatusInternalServerError)
		case utils.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(stringToJSON(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(stringToJSON("Checked out"))
}

// ImageRecognitionHandler uploads an image and use recognition software
func ImageRecognitionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("plate")
	if err != nil {
		logrus.Warn("Error getting plate image")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(stringToJSON(utils.ErrBadRequest.Error()))

		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("assets", "upload-*.png")
	if err != nil {
		logrus.Warn("Error creating a temp file")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(stringToJSON(utils.ErrInternalServer.Error()))

		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Warn("Error reading file")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(stringToJSON(utils.ErrBadRequest.Error()))

		return
	}

	tempFile.Write(fileBytes)

	request, err := usecases.Recognize(tempFile.Name())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(stringToJSON(utils.ErrImageRecognition.Error()))

		return
	}

	id, err := usecases.MakeReservation(request)
	if err != nil {
		switch err {
		case utils.ErrInternalServer:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(stringToJSON(err.Error()))
		case utils.ErrPlateNotValid:
			w.WriteHeader(http.StatusBadRequest)
			w.Write(stringToJSON(fmt.Sprintf("Text recognized `%s` is not in the right format: AAA-1234", request.Plate)))
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(idToJSON(id))
}

func stringToJSON(str string) []byte {
	res := struct {
		Response string `json:"response"`
	}{
		Response: str,
	}

	bytes, _ := json.Marshal(res)
	return bytes
}

func idToJSON(id uint) []byte {
	res := struct {
		Response uint `json:"id"`
	}{
		Response: id,
	}

	bytes, _ := json.Marshal(res)
	return bytes
}
