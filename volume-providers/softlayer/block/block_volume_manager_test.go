/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package softlayer_block

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend/fakes"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/common"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SLBlockSession", func() {
	var (
		fakeBakend *fakes.Session
		//locationSvc fakes.LocationService
		logger                     *zap.Logger
		datacenters                []datatypes.Location
		slSession                  *SLBlockSession
		volumeRequest              provider.Volume
		packageService             *fakes.ProductPackageService
		accountService             *fakes.AccountService
		networkStorageService      *fakes.NetworkStorageService
		networkStorageIscsiService *fakes.NetworkStorageIscsiService
		billingItemService         *fakes.BillingItemService
		enduranceType              = provider.VolumeProviderType("endurance")
		performaceType             = provider.VolumeProviderType("performance")
	)

	BeforeEach(func() {
		fakeBakend = fakes.NewSession()
		logger = zap.NewNop()
		config, _ := config.ReadConfig("", logger)
		slSession = &SLBlockSession{
			common.SLSession{
				Config:     config.Softlayer,
				Backend:    fakeBakend,
				Logger:     logger,
				VolumeType: "ISCSI",
			},
		}
	})

	Describe("VolumeCreate", func() {
		var (
			capacityItemPrices   []datatypes.Product_Item_Price
			itemPrices           map[string]*datatypes.Product_Item_Price
			storagespacePrice    []datatypes.Product_Item_Price
			itemPriceErrors      map[string]error
			productPackages      []datatypes.Product_Package
			networkStorages      []datatypes.Network_Storage
			networkIscsiStorages []datatypes.Network_Storage_Iscsi
		)

		BeforeEach(func() {
			capacityItemPrices = []datatypes.Product_Item_Price{
				{
					Id:                         sl.Int(55555),
					CapacityRestrictionType:    sl.String("STORAGE_TIER_LEVEL"),
					CapacityRestrictionMinimum: sl.String("100"),
					CapacityRestrictionMaximum: sl.String("100"),
					Item: &datatypes.Product_Item{
						Id:       sl.Int(55556),
						Capacity: sl.Float(250),
						KeyName:  sl.String("250_GB_PERFORMANCE_STORAGE_SPACE"),
						Units:    sl.String("GB"),
					},
				},
				{
					Id:                         sl.Int(33333),
					CapacityRestrictionMinimum: sl.String("100"),
					CapacityRestrictionMaximum: sl.String("100"),
					Item: &datatypes.Product_Item{
						Id:       sl.Int(33334),
						Capacity: sl.Float(20),
						KeyName:  sl.String("20_GB_PERFORMANCE_STORAGE_SPACE"),
						Units:    sl.String("GB"),
					},
				},
				{
					Id:                         sl.Int(44444),
					CapacityRestrictionMinimum: sl.String("100"),
					CapacityRestrictionMaximum: sl.String("100"),
					Item: &datatypes.Product_Item{
						Id:       sl.Int(44445),
						Capacity: sl.Float(100),
						KeyName:  sl.String("100_GB_PERFORMANCE_STORAGE_SPACE"),
						Units:    sl.String("GB"),
					},
				},
			}

			itemPrices = map[string]*datatypes.Product_Item_Price{
				"storage_as_a_service/STORAGE_AS_A_SERVICE": &datatypes.Product_Item_Price{
					Id: sl.Int(1000),
				},
				"storage_block/BLOCK_STORAGE_2": &datatypes.Product_Item_Price{
					Id: sl.Int(2000),
				},
				"storage_tier_level/LOW_INTENSITY_TIER": &datatypes.Product_Item_Price{
					Id: sl.Int(3000),
					Item: &datatypes.Product_Item{
						Id: sl.Int(3001),
						Attributes: []datatypes.Product_Item_Attribute{
							{AttributeTypeKeyName: sl.String("STORAGE_TIER_LEVEL"), Value: sl.String("100")},
						},
					},
				},
				"storage_tier_level/READHEAVY_TIER": &datatypes.Product_Item_Price{
					Id: sl.Int(4000),
					Item: &datatypes.Product_Item{
						Id: sl.Int(4001),
						Attributes: []datatypes.Product_Item_Attribute{
							{AttributeTypeKeyName: sl.String("STORAGE_TIER_LEVEL"), Value: sl.String("200")},
						},
					},
				},
				"storage_tier_level/WRITEHEAVY_TIER": &datatypes.Product_Item_Price{
					Id: sl.Int(5000),
					Item: &datatypes.Product_Item{
						Id: sl.Int(5001),
						Attributes: []datatypes.Product_Item_Attribute{
							{AttributeTypeKeyName: sl.String("STORAGE_TIER_LEVEL"), Value: sl.String("300")},
						},
					},
				},
				"storage_tier_level/10_IOPS_PER_GB": &datatypes.Product_Item_Price{
					Id: sl.Int(6000),
					Item: &datatypes.Product_Item{
						Id: sl.Int(6001),
						Attributes: []datatypes.Product_Item_Attribute{
							{AttributeTypeKeyName: sl.String("STORAGE_TIER_LEVEL"), Value: sl.String("1000")},
						},
					},
				},
				"storage_capacity": &datatypes.Product_Item_Price{
					Id: sl.Int(7000),
					Item: &datatypes.Product_Item{
						CapacityMaximum: sl.String("12000"), CapacityMinimum: sl.String("1"),
						Id: sl.Int(96), KeyName: sl.String("STORAGE_SPACE_FOR_2_IOPS_PER_GB"),
						Capacity: sl.Float(0),
					},
				},
			}
			storagespacePrice = []datatypes.Product_Item_Price{
				{
					Id: sl.Int(7000), Item: &datatypes.Product_Item{
						CapacityMaximum: sl.String("12000"), CapacityMinimum: sl.String("1"),
						Id: sl.Int(96), KeyName: sl.String("STORAGE_SPACE_FOR_2_IOPS_PER_GB"),
						Capacity: sl.Float(0),
					},
				},
				{
					Id: sl.Int(19001), Item: &datatypes.Product_Item{
						CapacityMaximum: sl.String("12000"), CapacityMinimum: sl.String("1000"),
						Id: sl.Int(9625), KeyName: sl.String("STORAGE_SPACE_FOR_10_IOPS_PER_GB"),
					},
				},
			}
			itemPriceErrors = map[string]error{}

			datacenters = []datatypes.Location{
				{Id: sl.Int(1), Name: sl.String("TEST-DC01")},
				{Id: sl.Int(2), Name: sl.String("TEST-DC02")},
			}

			fakeBakend.GetLocationServiceFake().GetDatacentersReturns(datacenters, nil)
			productPackages = []datatypes.Product_Package{
				{
					Id:         sl.Int(1000),
					Name:       sl.String("storage_as_a_service"),
					ItemPrices: storagespacePrice,
					Categories: []datatypes.Product_Item_Category{
						{CategoryCode: sl.String("storage_as_a_service")},
					},
				},
			}

			networkStorages = []datatypes.Network_Storage{
				{
					Id: sl.Int(1111),
					BillingItem: &datatypes.Billing_Item{
						Id: sl.Int(2222),
						OrderItem: &datatypes.Billing_Order_Item{
							Order: &datatypes.Billing_Order{
								Id: sl.Int(1111),
							},
						},
					},
					ActiveTransactionCount: sl.Uint(0),
					StorageType: &datatypes.Network_Storage_Type{
						KeyName: sl.String("endurance"),
					},
				},
				{
					Id: sl.Int(2222),
					BillingItem: &datatypes.Billing_Item{
						Id: sl.Int(2222),
						OrderItem: &datatypes.Billing_Order_Item{
							Order: &datatypes.Billing_Order{
								Id: sl.Int(2222),
							},
						},
					},
					ActiveTransactionCount: sl.Uint(0),
					StorageType: &datatypes.Network_Storage_Type{
						KeyName: sl.String("performance"),
					},
				},
			}

			networkIscsiStorages = []datatypes.Network_Storage_Iscsi{
				{Network_Storage: networkStorages[0]},
				{Network_Storage: networkStorages[1]},
			}

			volumeRequest = provider.Volume{
				VolumeType:      "block",
				ProviderType:    enduranceType,
				Capacity:        sl.Int(20),
				Tier:            sl.String("0.25"),
				ServiceOffering: sl.String("storage_as_a_service"),
				Az:              "TEST-DC01",
			}
			packageService = fakeBakend.GetProductPackageServiceFake()
			packageService.FilterReturns(packageService)
			packageService.MaskReturns(packageService)
			packageService.GetAllObjectsReturns(productPackages, nil)

			accountService = fakeBakend.GetAccountServiceFake()
			accountService.FilterReturns(accountService)
			accountService.MaskReturns(accountService)
			//accountService.GetNetworkStorageReturns(networkStorages, nil)
			accountService.GetNetworkStorageCalls(func() ([]datatypes.Network_Storage, error) {
				if volumeRequest.ProviderType == enduranceType {
					return networkStorages[0:1], nil
				}
				if volumeRequest.ProviderType == performaceType {
					return networkStorages[1:2], nil
				}
				return networkStorages[0:2], nil

			})

			networkStorageService = fakeBakend.GetNetworkStorageServiceFake()
			networkStorageService.FilterReturns(networkStorageService)
			networkStorageService.MaskReturns(networkStorageService)
			networkStorageService.IDReturns(networkStorageService)
			//networkStorageService.GetObjectReturns(networkStorages[0], nil)
			networkStorageService.GetObjectCalls(func() (datatypes.Network_Storage, error) {
				if volumeRequest.ProviderType == enduranceType {
					return networkStorages[0], nil
				} else {
					return networkStorages[1], nil
				}
			})

			networkStorageIscsiService = fakeBakend.GetNetworkStorageIscsiServiceFake()
			networkStorageIscsiService.FilterReturns(networkStorageIscsiService)
			networkStorageIscsiService.MaskReturns(networkStorageIscsiService)
			networkStorageIscsiService.IDReturns(networkStorageIscsiService)
			//networkStorageIscsiService.GetObjectReturns(networkIscsiStorages[0], nil)

			billingItemService = fakeBakend.GetBillingItemServiceFake()
			billingItemService.FilterReturns(billingItemService)
			billingItemService.MaskReturns(billingItemService)
			billingItemService.IDReturns(billingItemService)
			billingItemService.CancelItemReturns(true, nil)

			networkStorageIscsiService.GetObjectCalls(func() (datatypes.Network_Storage_Iscsi, error) {
				if volumeRequest.ProviderType == enduranceType {
					return networkIscsiStorages[0], nil
				} else {
					return networkIscsiStorages[1], nil
				}
			})
			productOrderReceipt := datatypes.Container_Product_Order_Receipt{
				OrderId: sl.Int(1111),
			}
			fakeBakend.GetProductOrderServiceFake().PlaceOrderReturns(productOrderReceipt, nil)
		})

		It("orders endurance storage", func() {
			volume, err := slSession.VolumeCreate(volumeRequest)

			Expect(err).NotTo(HaveOccurred())
			Expect(volume.VolumeID).To(Equal("1111"))
			Expect(volume.ProviderType).To(Equal(provider.VolumeProviderType("endurance")))
			Expect(accountService.GetNetworkStorageCallCount()).To(Equal(1))
			Expect(packageService.GetAllObjectsCallCount()).To(Equal(1))
			Expect(networkStorageService.GetObjectCallCount()).To(Equal(1))
			Expect(networkStorageIscsiService.GetObjectCallCount()).To(Equal(1))

		})
		It("orders performance storage", func() {
			volumeRequest.ProviderType = performaceType
			volumeRequest.Iops = sl.String("100")
			volume, err := slSession.VolumeCreate(volumeRequest)

			Expect(err).NotTo(HaveOccurred())
			Expect(volume.VolumeID).To(Equal("2222"))
			Expect(volume.ProviderType).To(Equal(provider.VolumeProviderType("performance")))
			Expect(accountService.GetNetworkStorageCallCount()).To(Equal(1))
			Expect(packageService.GetAllObjectsCallCount()).To(Equal(1))
			Expect(networkStorageService.GetObjectCallCount()).To(Equal(1))
			Expect(networkStorageIscsiService.GetObjectCallCount()).To(Equal(1))

		})
		Context("when ordering fails with wrong provider type ", func() {
			It("returns an error", func() {
				volumeRequest.ProviderType = "WRONG"
				_, err := slSession.VolumeCreate(volumeRequest)
				Expect(err.Error()).To(ContainSubstring("Storage type is wrong or not provided"))
			})
		})

		It("delete volume success", func() {
			volume, err := slSession.VolumeGet("2222")
			err = slSession.VolumeDelete(volume)
			Expect(err).NotTo(HaveOccurred())
			Expect(networkStorageIscsiService.GetObjectCallCount()).To(Equal(2))
			Expect(billingItemService.CancelItemCallCount()).To(Equal(1))

		})

		It("get volume  by request success", func() {
			volume, err := slSession.GetVolumeByRequestID("1111")
			Expect(err).NotTo(HaveOccurred())
			Expect(volume.VolumeID).To(Equal("1111"))
			Expect(accountService.GetNetworkStorageCallCount()).To(Equal(1))
		})

		It("get volume  by request fail", func() {
			_, err := slSession.GetVolumeByRequestID("2222")
			Expect(err.Error()).To(ContainSubstring("Incorrect storage found"))
		})

	})
})
