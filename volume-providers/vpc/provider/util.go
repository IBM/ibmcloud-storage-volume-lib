/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// maxRetryAttempt ...
var maxRetryAttempt = 10

// maxRetryGap ...
var maxRetryGap = 60

// retryGap ...
var retryGap = 10

var volumeIDPartsCount = 5

var skipErrorCodes = map[string]bool{
	"validation_invalid_name":          true,
	"volume_capacity_max":              true,
	"volume_id_invalid":                true,
	"volume_profile_iops_invalid":      true,
	"volume_capacity_zero_or_negative": true,
	"not_found":                        true,
	"internal_error":                   false,
	"invalid_route":                    false,

	// IKS ms error code for skip re-try
	"ST0008": true, //resources not found
	"ST0005": true, //worker node could not be found
}

// retry ...
func retry(logger *zap.Logger, retryfunc func() error) error {
	var err error

	for i := 0; i < maxRetryAttempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(retryGap) * time.Second)
		}
		err = retryfunc()
		if err != nil {
			//Skip retry for the below type of Errors
			modelError, ok := err.(*models.Error)
			if !ok {
				continue
			}
			if skipRetry(modelError) {
				break
			}
			if i >= 1 {
				retryGap = 2 * retryGap
				if retryGap > maxRetryGap {
					retryGap = maxRetryGap
				}
			}
			if (i + 1) < maxRetryAttempt {
				logger.Info("Error while executing the function. Re-attempting execution ..", zap.Int("attempt..", i+2), zap.Int("retry-gap", retryGap), zap.Int("max-retry-Attempts", maxRetryGap), zap.Error(err))
			}
			continue
		}
		return err
	}
	return err
}

// skipRetry skip retry as per listed error codes
func skipRetry(err *models.Error) bool {
	for _, errorItem := range err.Errors {
		skipStatus, ok := skipErrorCodes[string(errorItem.Code)]
		if ok {
			return skipStatus
		}
	}
	return false
}

// skipRetryForIKS skip retry as per listed error codes
func skipRetryForIKS(err *models.IksError) bool {
	skipStatus, ok := skipErrorCodes[string(err.Code)]
	if ok {
		return skipStatus
	}
	return false
}

// skipRetryForAttach skip retry as per listed error codes
func skipRetryForAttach(err error, isIKS bool) bool {
	// Only for storage-api ms related calls error
	if isIKS {
		iksError, ok := err.(*models.IksError)
		if ok {
			return skipRetryForIKS(iksError)
		}
		return false
	}

	// Only for RIaaS attachment related calls error
	riaasError, ok := err.(*models.Error)
	if ok {
		return skipRetry(riaasError)
	}
	return false
}

// FlexyRetry ...
type FlexyRetry struct {
	maxRetryAttempt int
	maxRetryGap     int
}

// NewFlexyRetryDefault ...
func NewFlexyRetryDefault() FlexyRetry {
	return FlexyRetry{
		// Default values as we configuration
		maxRetryAttempt: maxRetryAttempt,
		maxRetryGap:     maxRetryGap,
	}
}

// NewFlexyRetry ...
func NewFlexyRetry(maxRtyAtmpt int, maxrRtyGap int) FlexyRetry {
	return FlexyRetry{
		maxRetryAttempt: maxRtyAtmpt,
		maxRetryGap:     maxrRtyGap,
	}
}

// FlexyRetry ...
func (fRetry *FlexyRetry) FlexyRetry(logger *zap.Logger, funcToRetry func() (error, bool)) error {
	var err error
	var stopRetry bool
	for i := 0; i < fRetry.maxRetryAttempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(retryGap) * time.Second)
		}
		// Call function which required retry, retry is decided by funtion itself
		err, stopRetry = funcToRetry()
		if stopRetry {
			break
		}

		// Update retry gap as per exponentioal
		if i >= 1 {
			retryGap = 2 * retryGap
			if retryGap > fRetry.maxRetryGap {
				retryGap = fRetry.maxRetryGap
			}
		}
		if (i + 1) < fRetry.maxRetryAttempt {
			logger.Info("UNEXPECTED RESULT, Re-attempting execution ..", zap.Int("attempt..", i+2),
				zap.Int("retry-gap", retryGap), zap.Int("max-retry-Attempts", fRetry.maxRetryAttempt),
				zap.Bool("stopRetry", stopRetry), zap.Error(err))
		}
	}
	return err
}

// FlexyRetryWithConstGap ...
func (fRetry *FlexyRetry) FlexyRetryWithConstGap(logger *zap.Logger, funcToRetry func() (error, bool)) error {
	var err error
	var stopRetry bool
	// lets have more number of try for wait for attach and detach specially
	totalAttempt := fRetry.maxRetryAttempt * 4 // 40 time as per default values i.e 400 seconds
	for i := 0; i < totalAttempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(retryGap) * time.Second)
		}
		// Call function which required retry, retry is decided by funtion itself
		err, stopRetry = funcToRetry()
		if stopRetry {
			break
		}

		if (i + 1) < totalAttempt {
			logger.Info("UNEXPECTED RESULT from FlexyRetryWithConstGap, Re-attempting execution ..", zap.Int("attempt..", i+2),
				zap.Int("retry-gap", retryGap), zap.Int("max-retry-Attempts", totalAttempt),
				zap.Bool("stopRetry", stopRetry), zap.Error(err))
		}
	}
	return err
}

// ToInt ...
func ToInt(valueInInt string) int {
	value, err := strconv.Atoi(valueInInt)
	if err != nil {
		return 0
	}
	return value
}

// ToInt64 ...
func ToInt64(valueInInt string) int64 {
	value, err := strconv.ParseInt(valueInInt, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// FromProviderToLibVolume converting vpc provider volume type to generic lib volume type
func FromProviderToLibVolume(vpcVolume *models.Volume, logger *zap.Logger) (libVolume *provider.Volume) {
	logger.Debug("Entry of FromProviderToLibVolume method...")
	defer logger.Debug("Exit from FromProviderToLibVolume method...")

	if vpcVolume == nil {
		logger.Info("Volume details are empty")
		return
	}

	if vpcVolume.Zone == nil {
		logger.Info("Volume zone is empty")
		return
	}

	logger.Debug("Volume details of VPC client", zap.Reflect("models.Volume", vpcVolume))

	volumeCap := int(vpcVolume.Capacity)
	iops := strconv.Itoa(int(vpcVolume.Iops))
	var createdDate time.Time
	if vpcVolume.CreatedAt != nil {
		createdDate = *vpcVolume.CreatedAt
	}

	libVolume = &provider.Volume{
		VolumeID:     vpcVolume.ID,
		Provider:     VPC,
		Capacity:     &volumeCap,
		Iops:         &iops,
		VolumeType:   VolumeType,
		CreationTime: createdDate,
	}
	if vpcVolume.Zone != nil {
		libVolume.Region = vpcVolume.Zone.Name
	}
	return
}

// IsValidVolumeIDFormat validating
func IsValidVolumeIDFormat(volID string) bool {
	parts := strings.Split(volID, "-")
	if len(parts) != volumeIDPartsCount {
		return false
	}
	return true
}

// SetRetryParameters sets the retry logic parameters
func SetRetryParameters(maxAttempts int, maxGap int) {
	if maxAttempts > 0 {
		maxRetryAttempt = maxAttempts
	}

	if maxGap > 0 {
		maxRetryGap = maxGap
	}
}
