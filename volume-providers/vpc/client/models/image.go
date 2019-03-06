/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package models

import "time"

// ImageVisibilityType ...
type ImageVisibilityType string

// String ...
func (i ImageVisibilityType) String() string { return string(i) }

// Image visibility values
const (
	ImagePublic  ImageVisibilityType = "public"
	ImagePrivate ImageVisibilityType = "private"
)

// ImageFormatType ...
type ImageFormatType string

// String ...
func (i ImageFormatType) String() string { return string(i) }

// Image format values
const (
	ImageFormatBox   ImageFormatType = "box"
	ImageFormatOVA   ImageFormatType = "ova"
	ImageFormatQCOW2 ImageFormatType = "qcow2"
	ImageFormatRaw   ImageFormatType = "raw"
	ImageFormatVDI   ImageFormatType = "vdi"
	ImageFormatVHDX  ImageFormatType = "vhdx"
	ImageFormatVMDK  ImageFormatType = "vmdk"
)

// ImageStatusType ...
type ImageStatusType string

// String ...
func (i ImageStatusType) String() string { return string(i) }

// Image status values
const (
	ImageStatusAvailable ImageFormatType = "available"
	ImageStatusFailed    ImageFormatType = "failed"
	ImageStatusPending   ImageFormatType = "pending"
)

// Image ...
type Image struct {
	Architecture    string              `json:"architecture,omitempty"`
	CreatedAt       *time.Time          `json:"created_at,omitempty"`
	CRN             string              `json:"crn,omitempty"`
	Description     string              `json:"description,omitempty"`
	File            *ImageFile          `json:"file,omitempty"`
	Format          ImageFormatType     `json:"format,omitempty"`
	ID              string              `json:"id,omitempty"`
	Name            string              `json:"name,omitempty"`
	OperatingSystem *OperatingSystem    `json:"operating_system,omitempty"`
	ResourceGroup   *ResourceGroup      `json:"resource_group,omitempty"`
	Status          ImageStatusType     `json:"status,omitempty"`
	Tags            []string            `json:"tags,omitempty"`
	Visibility      ImageVisibilityType `json:"visibility,omitempty"`
}
