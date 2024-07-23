package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"time"
)

var (
	ssoExample = `  # Login using a single sign-on provider
  omnistrate-ctl sso`
)

// SsoCmd represents the login command
var SsoCmd = &cobra.Command{
	Use:          `sso`,
	Short:        "Log in using a single sign-on provider.",
	Long:         `The sso command is used to authenticate and log in to the Omnistrate platform.`,
	Example:      ssoExample,
	RunE:         run,
	SilenceUsage: true,
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
	// Copy the user code to the clipboard
	err = clipboard.WriteAll(deviceCodeResponse.UserCode)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error copying user code to clipboard: %v\n", err))
		utils.PrintError(err)
		return err
	}

	// Automatically open the verification URI in the default browser
	fmt.Println("Attempting to automatically open the SSO authentication page in your default browser.")
	err = browser.OpenURL(deviceCodeResponse.VerificationURI)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error opening browser: %v\n", err))
		utils.PrintError(err)
		return err
	}
	fmt.Print("If the browser does not open or you wish to use a different device to authorize this request, open the following URL:\n\n")
	fmt.Printf("%s\n\n", deviceCodeResponse.VerificationURI)
	fmt.Print("The code has been copied to your clipboard. Paste it in the browser when prompted.\n")
	fmt.Print("You can also manually type in the code:\n\n")
	fmt.Printf("%s\n\n", deviceCodeResponse.UserCode)

	// Step 3: Poll GitHub to check if the user authorized the device
	_, err = pollForAccessToken(deviceCodeResponse.DeviceCode, deviceCodeResponse.Interval)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error polling for access token: %v\n", err))
		utils.PrintError(err)
		return err
	}

	// Step 4: Use the access token to authenticate with the Omnistrate platform
	request := signinapi.LoginWithIdentityProviderRequest{
		AuthorizationCode:    deviceCodeResponse.DeviceCode,
		IdentityProviderName: signinapi.IdentityProviderName("GitHub"),
		RedirectURI:          commonutils.ToPtr("https://omnistrate.dev/idp-auth"),
	}

	res, err := dataaccess.LoginWithIdentityProvider(request)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error logging in with identity provider: %v\n", err))
		utils.PrintError(err)
		return err
	}

	token := res.JWTToken

	authConfig := config.AuthConfig{
		Email: email,
		Token: token,
	}
	if err = config.CreateOrUpdateAuthConfig(authConfig); err != nil {
		utils.PrintError(err)
		return err
	}

	utils.PrintSuccess("Successfully logged in.")

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
	devClientID  = "Ov23ctpQGrpGvsIIJxFv"
	prodClientID = "Ov23li2nyhdelepEtjcg"
	scope        = "user"
)

// requestDeviceCode requests a device and user verification code from GitHub
func requestDeviceCode() (*DeviceCodeResponse, error) {
	clientID := devClientID
	if utils.IsProd() {
		clientID = prodClientID
	}
	data := map[string]string{
		"client_id": clientID,
		"scope":     scope,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", "https://github.com/login/device/code", bytes.NewBuffer(dataBytes))
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

// pollForAccessToken polls GitHub for an access token
func pollForAccessToken(deviceCode string, interval int) (*AccessTokenResponse, error) {
	clientID := devClientID
	if utils.IsProd() {
		clientID = prodClientID
	}
	data := map[string]string{
		"client_id":   clientID,
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(context.Background(), "POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(dataBytes))
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

		//println(accessTokenResponse.AccessToken)
		//println(accessTokenResponse.TokenType)
		//println(accessTokenResponse.Scope)
		return &accessTokenResponse, nil
	}
}
