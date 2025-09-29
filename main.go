package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	
	sp "github.com/amzapi/selling-partner-api-sdk/pkg/selling-partner"
	"github.com/amzapi/selling-partner-api-sdk/sellers"

	"github.com/joho/godotenv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}
	
	sellingPartner, err := sp.NewSellingPartner(&sp.Config{
		ClientID:     os.Getenv("SP_API_CLIENT_ID"),
		ClientSecret: os.Getenv("SP_API_CLIENT_SECRET"),
		RefreshToken: os.Getenv("SP_API_REFRESH_TOKEN"),
	})

	if err != nil {
		panic(err)
	}

	endpoint := "https://sellingpartnerapi-na.amazon.com"

	seller, err := sellers.NewClientWithResponses(endpoint,
		sellers.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("X-Amzn-Requestid", uuid.New().String()) //tracking requests
			err = sellingPartner.AuthorizeRequest(req)
			if err != nil {
				return errors.Wrap(err, "sign error")
			}
			dump, err := httputil.DumpRequest(req, true)
			if err != nil {
				return errors.Wrap(err, "DumpRequest Error")
			}
			log.Printf("DumpRequest = %s", dump)
			return nil
		}),
		sellers.WithResponseAfter(func(ctx context.Context, rsp *http.Response) error {
			dump, err := httputil.DumpResponse(rsp, true)
			if err != nil {
				return errors.Wrap(err, "DumpResponse Error")
			}
			log.Printf("DumpResponse = %s", dump)
			return nil
		}),
	)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	_, err = seller.GetMarketplaceParticipationsWithResponse(ctx)

	if err != nil {
		panic(err)
	}
}
