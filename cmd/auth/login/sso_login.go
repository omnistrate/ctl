package login

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func SSOLogin(identityProviderName string) error {
	// Step 1: Request device and user verification codes
	deviceCodeResponse, err := requestDeviceCode(identityProviderName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error requesting device code: %v\n", err))
		utils.PrintError(err)
		return err
	}

	// Step 2: Prompt the user to enter the user code in a browser
	// Copy the user code to the clipboard
	err = clipboard.WriteAll(deviceCodeResponse.UserCode)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error copying user code to clipboard: %v\n", err))
		utils.PrintError(err)
		return err
	}

	// Automatically open the verification URI in the default browser
	fmt.Println("Attempting to automatically open the SSO authentication page in your default browser.")
	err = browser.OpenURL(getVerificationURI(identityProviderName))
	if err != nil {
		err = errors.New(fmt.Sprintf("Error opening browser: %v\n", err))
		utils.PrintError(err)
		return err
	}
	fmt.Print("If the browser does not open or you wish to use a different device to authorize this request, open the following URL:\n\n")
	fmt.Printf("%s\n\n", getVerificationURI(identityProviderName))
	fmt.Print("The code has been copied to your clipboard. Paste it in the browser when prompted.\n")
	fmt.Print("You can also manually type in the code:\n\n")
	fmt.Printf("%s\n\n", deviceCodeResponse.UserCode)

	// Step 3: Poll identity provider server to check if the user authorized the device via backend API
	jwtTokenResponse, err := pollForAccessTokenAndLogin(identityProviderName, deviceCodeResponse.DeviceCode, deviceCodeResponse.Interval)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error polling for access token: %v\n", err))
		utils.PrintError(err)
		return err
	}

	token := jwtTokenResponse.JWTToken

	authConfig := config.AuthConfig{
		Token: token,
	}
	if err = config.CreateOrUpdateAuthConfig(authConfig); err != nil {
		utils.PrintError(err)
		return err
	}

	utils.PrintSuccess("Successfully logged in")

	return nil
}

// DeviceCodeResponse represents the response from the device code request
type DeviceCodeResponse struct {
	DeviceCode string `json:"device_code"`
	UserCode   string `json:"user_code"`
	ExpiresIn  int    `json:"expires_in"`
	Interval   int    `json:"interval"`
}

// AccessTokenResponse represents the response from the jwt token request
type AccessTokenResponse struct {
	JWTToken string `json:"jwt_token"`
}

// GitHub client credentials
const (
	gitHubDevClientID  = "Ov23ctpQGrpGvsIIJxFv"
	gitHubProdClientID = "Ov23li2nyhdelepEtjcg"
	googleDevClientID  = "635031719937-gqvm0qeelipdc812g9ie2v6ohk3j6gs6.apps.googleusercontent.com" // #nosec G101
	googleProdClientID = "421577562987-98lkfnu7e07rig5p6rt4p0dgqpktihhb.apps.googleusercontent.com" // #nosec G101
	gitHubScope        = "user:email"
	googleScope        = "email"
)

// requestDeviceCode requests a device and user verification code from GitHub
func requestDeviceCode(identityProviderName string) (*DeviceCodeResponse, error) {
	data := map[string]string{
		"client_id": getClientID(identityProviderName),
		"scope":     getScope(identityProviderName),
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", getDeviceCodeURL(identityProviderName), bytes.NewBuffer(dataBytes))
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

	body, err := io.ReadAll(resp.Body)
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

// pollForAccessTokenAndLogin polls identity provider server for an access token and uses it to log user into the platform
func pollForAccessTokenAndLogin(identityProviderName, deviceCode string, interval int) (*AccessTokenResponse, error) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		resp, err := dataaccess.LoginWithIdentityProvider(deviceCode, identityProviderName)
		if err != nil {
			if strings.Contains(err.Error(), "Failed to get access token with status code: 428 Precondition Required") { // authorization_pending
				continue
			}
			if strings.Contains(err.Error(), "Failed to get access token with status code: 403 Forbidden") { // access_denied
				return nil, errors.New("Access denied. Please try again.")
			}
			if identityProviderName == "GitHub" && strings.Contains(err.Error(), "Invalid request: empty access token") { // authorization_pending
				continue
			}
			return nil, err
		}

		return &AccessTokenResponse{
			JWTToken: resp.JWTToken,
		}, nil
	}
}

func getClientID(identityProviderName string) string {
	if identityProviderName == "Google for CTL" {
		clientID := googleDevClientID
		if utils.IsProd() {
			clientID = googleProdClientID
		}
		return clientID
	}

	// Default to GitHub
	clientID := gitHubDevClientID
	if utils.IsProd() {
		clientID = gitHubProdClientID
	}
	return clientID
}

func getScope(identityProviderName string) string {
	if identityProviderName == "Google for CTL" {
		return googleScope
	}

	// Default to GitHub
	return gitHubScope
}

func getDeviceCodeURL(identityProviderName string) string {
	if identityProviderName == "Google for CTL" {
		return "https://oauth2.googleapis.com/device/code"
	}

	// Default to GitHub
	return "https://github.com/login/device/code"
}

func getVerificationURI(identityProviderName string) string {
	if identityProviderName == "Google for CTL" {
		return "https://www.google.com/device"
	}

	// Default to GitHub
	return "https://github.com/login/device"
}
