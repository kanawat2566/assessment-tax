package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	cm "github.com/kanawat2566/assessment-tax/common"
	ct "github.com/kanawat2566/assessment-tax/constants"
	md "github.com/kanawat2566/assessment-tax/model"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/labstack/echo/v4"
)

type taxHandler struct {
	serv services.TaxService
}

type CustomValidator struct {
	Validator *validator.Validate
}

func NewHandler(s services.TaxService) *taxHandler {
	return &taxHandler{serv: s}
}

func (h *taxHandler) CalculationsHandler(c echo.Context) error {
	rq := new(md.TaxRequest)

	if err := BindWithValidate(c, rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.serv.TaxCalculations(*rq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if res.TaxRefund > 0 {
		return c.JSON(http.StatusOK, res)
	} else {

		tax := md.TaxLevelReponse{Tax: res.Tax, TaxLevels: res.TaxLevels}
		return c.JSON(http.StatusOK, tax)
	}

}

func (h *taxHandler) Deductions(c echo.Context) error {
	rq := new(md.DeductRequest)
	d := c.Param("type")

	dd := ct.Deductions[d]
	if len(dd.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ct.ErrMsgDeductNotFound)

	}

	if err := BindWithValidate(c, rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.serv.SetAdminDeductions(ct.Deduction{Type: dd.Type, Amount: rq.Amount})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := map[string]interface{}{
		res.Name: res.Amount,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *taxHandler) CalFromUploadCsvHandler(c echo.Context) error {

	uploadType := c.Param("uploadType")

	if uploadType != ct.PathParamUploadCsv {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New(ct.ErrMsgInvalidPathParam))
	}

	csv, err := UploadFromCsv(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	taxes, err := h.serv.TaxCalFromCsv(csv)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"taxes": taxes,
	}

	return c.JSON(http.StatusOK, response)

}

func UploadFromCsv(c echo.Context) ([]md.TaxRequest, error) {

	file, err := c.FormFile("taxFile")
	if err != nil {
		return nil, errors.New(ct.ErrMsgFileNoUpload)
	}
	src, err := file.Open()
	if err != nil {
		return nil, errors.New(ct.ErrMsgReadCsvFailed)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, errors.New(ct.ErrMsgReadCsvFailed)
	}

	reader := csv.NewReader(bytes.NewReader(fileBytes))

	header, err := reader.Read()
	if !reflect.DeepEqual(header, ct.CsvFomatFile) {
		return nil, errors.New(ct.ErrMsgCsvInvaildFormat)
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New(ct.ErrMsgCsvInvaildFormat)
	}

	var taxReqs []md.TaxRequest
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		taxReq := md.TaxRequest{}
		if len(row) != 3 {
			return nil, errors.New(ct.ErrMsgCsvInvaildFormat)
		}

		if taxReq.TotalIncome, err = strconv.ParseFloat(row[0], 64); err != nil {
			return nil, errors.New(cm.MsgWithInt(ct.ErrInvalidIncomeCsv, i+2))
		}
		if taxReq.WHT, err = strconv.ParseFloat(row[1], 64); err != nil {
			return nil, errors.New(cm.MsgWithInt(ct.ErrInvalidWHTCsv, i+2))
		}
		var donation float64
		if donation, err = strconv.ParseFloat(row[2], 64); err != nil {
			return nil, errors.New(cm.MsgWithInt(ct.ErrInvalidDonationCsv, i+2))
		}
		taxReq.Allowances = []md.Allowance{
			{AllowanceType: ct.Donation, Amount: donation},
		}

		taxReqs = append(taxReqs, taxReq)
	}

	return taxReqs, nil

}
