import (
	"context"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
)

func Handler(c echo.Context) error {
	ctx := c.Request().Context()

	client := &http.Client{Transport: &ochttp.Transport{}}
	req, err := http.NewRequest("GET", "http://service_a:8080/", nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
}
