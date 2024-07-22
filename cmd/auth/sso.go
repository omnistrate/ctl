package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ssoExample = `  # Login with github
  omnistrate-ctl sso github

  # Login with google
  omnistrate-ctl sso google`
)

// SsoCmd represents the login command
var SsoCmd = &cobra.Command{
	Use:          `sso [github|google]`,
	Short:        "Log in using a single sign-on provider.",
	Long:         `The sso command is used to authenticate and log in to the Omnistrate platform.`,
	Example:      ssoExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	SsoCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func run(cmd *cobra.Command, args []string) error {
	// Step 1: Request device and user verification codes from GitHub
	deviceCodeResponse, err := requestDeviceCode()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error requesting device code: %v\n", err))
		utils.PrintError(err)
		return err
	}

	// Step 2: Prompt the user to enter the user code in a browser
	fmt.Printf("Please enter the code %s at %s\n", deviceCodeResponse.UserCode, deviceCodeResponse.VerificationURI)

	// Step 3: Poll GitHub to check if the user authorized the device
	accessTokenResponse, err := pollForAccessToken(deviceCodeResponse.DeviceCode, deviceCodeResponse.Interval)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error polling for access token: %v\n", err))
		utils.PrintError(err)
		return err
	}

	fmt.Printf("Access token: %s\n", accessTokenResponse.AccessToken)

	return nil
}

// DeviceCodeResponse represents the response from the device code request
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// AccessTokenResponse represents the response from the access token request
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHub client credentials
const (
	devClientID = "Ov23ctpQGrpGvsIIJxFv"
)

// requestDeviceCode requests a device and user verification code from GitHub
func requestDeviceCode() (*DeviceCodeResponse, error) {
	data := map[string]string{
		"client_id": devClientID,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://github.com/login/device/code", bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var deviceCodeResponse DeviceCodeResponse
	err = json.Unmarshal(body, &deviceCodeResponse)
	if err != nil {
		return nil, err
	}

	return &deviceCodeResponse, nil
}

// pollForAccessToken polls GitHub for an access token
func pollForAccessToken(deviceCode string, interval int) (*AccessTokenResponse, error) {
	data := map[string]string{
		"client_id":   devClientID,
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(dataBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var errorResponse struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
			Interval         int    `json:"interval"`
		}
		if err = json.Unmarshal(body, &errorResponse); err == nil {
			if errorResponse.Error == "authorization_pending" {
				continue
			}
			if errorResponse.Error == "slow_down" {
				interval += 5
				continue
			}
		}

		var accessTokenResponse AccessTokenResponse
		err = json.Unmarshal(body, &accessTokenResponse)
		if err != nil {
			return nil, err
		}

		return &accessTokenResponse, nil
	}
}
