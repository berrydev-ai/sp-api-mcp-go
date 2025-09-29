# SP-API MCP Server Implementation Priorities

This document outlines all available Amazon Selling Partner API functions and their implementation priorities for the MCP server.

## Implementation Status Summary

**Currently Implemented:** 12 functions
- ✅ Authorization: GetAuthorizationCode  
- ✅ Orders: GetOrders, GetOrder, GetOrderAddress, GetOrderBuyerInfo, GetOrderItems, GetOrderItemsBuyerInfo
- ✅ Sales: GetOrderMetrics
- ✅ Reports: GetReports, CreateReport, GetReport, GetReportDocument

**Available for Implementation:** 100+ functions across 22 API modules

## Complete SP-API Function Inventory

### 1. **Authorization API** (`authorization`)
- ✅ **GetAuthorizationCode** - Get authorization code for delegated access

### 2. **Catalog Items API** (`catalog`)
- 🔲 **SearchCatalogItems** - Search for catalog items
- 🔲 **GetCatalogItem** - Get details for a specific catalog item

### 3. **FBA Inbound API** (`fbaInbound`)
- 🔲 **GetInboundGuidance** - Get inbound guidance for items
- 🔲 **CreateInboundShipmentPlan** - Create inbound shipment plan
- 🔲 **CreateInboundShipment** - Create inbound shipment
- 🔲 **UpdateInboundShipment** - Update inbound shipment
- 🔲 **GetPreorderInfo** - Get preorder information
- 🔲 **ConfirmPreorder** - Confirm preorder
- 🔲 **GetPrepInstructions** - Get prep instructions
- 🔲 **GetTransportDetails** - Get transport details
- 🔲 **PutTransportDetails** - Update transport details
- 🔲 **VoidTransport** - Void transport
- 🔲 **EstimateTransport** - Estimate transport
- 🔲 **ConfirmTransport** - Confirm transport
- 🔲 **GetLabels** - Get inbound shipment labels
- 🔲 **GetBillOfLading** - Get bill of lading
- 🔲 **GetShipments** - List inbound shipments
- 🔲 **GetShipmentItemsByShipmentId** - Get shipment items
- 🔲 **GetShipmentItems** - Get shipment items across shipments

### 4. **FBA Inventory API** (`fbaInventory`)
- 🔲 **GetInventorySummaries** - Get inventory summaries

### 5. **FBA Outbound API** (`fbaOutbound`)
- 🔲 **GetFulfillmentPreview** - Get fulfillment preview
- 🔲 **CreateFulfillmentOrder** - Create fulfillment order
- 🔲 **UpdateFulfillmentOrder** - Update fulfillment order
- 🔲 **CancelFulfillmentOrder** - Cancel fulfillment order
- 🔲 **GetFulfillmentOrder** - Get fulfillment order details
- 🔲 **ListAllFulfillmentOrders** - List all fulfillment orders
- 🔲 **GetPackageTrackingDetails** - Get package tracking details
- 🔲 **ListReturnReasonCodes** - List return reason codes
- 🔲 **CreateFulfillmentReturn** - Create fulfillment return
- 🔲 **GetFulfillmentReturn** - Get fulfillment return
- 🔲 **ListReturnReasonCodes** - List return reason codes
- 🔲 **GetFeatures** - Get available features
- 🔲 **GetFeatureInventory** - Get feature inventory
- 🔲 **GetFeatureSKU** - Get feature SKU

### 6. **Feeds API** (`feeds`)
- 🔲 **GetFeeds** - Get feed processing reports
- 🔲 **CreateFeed** - Create a feed
- 🔲 **GetFeed** - Get feed details
- 🔲 **CancelFeed** - Cancel a feed
- 🔲 **CreateFeedDocument** - Create feed document
- 🔲 **GetFeedDocument** - Get feed document

### 7. **Product Fees API** (`fees`)
- 🔲 **GetMyFeesEstimateForSKU** - Get fees estimate for SKU
- 🔲 **GetMyFeesEstimateForASIN** - Get fees estimate for ASIN

### 8. **Finances API** (`finances`)
- 🔲 **ListFinancialEventGroups** - List financial event groups
- 🔲 **ListFinancialEventsByGroupId** - List financial events by group
- 🔲 **ListFinancialEventsByOrderId** - List financial events by order
- 🔲 **ListFinancialEvents** - List all financial events

### 9. **Listings Items API** (`listingsItems`)
- 🔲 **DeleteListingsItem** - Delete a listing
- 🔲 **GetListingsItem** - Get listing details
- 🔲 **PutListingsItem** - Create or update listing
- 🔲 **PatchListingsItem** - Partially update listing

### 10. **Merchant Fulfillment API** (`merchantFulfillment`)
- 🔲 **GetEligibleShipmentServices** - Get eligible shipment services
- 🔲 **GetShipment** - Get shipment details
- 🔲 **CancelShipment** - Cancel shipment
- 🔲 **CancelShipmentOld** - Cancel shipment (legacy)
- 🔲 **CreateShipment** - Create shipment
- 🔲 **GetAdditionalSellerInputs** - Get additional seller inputs
- 🔲 **GetAdditionalSellerInputsOld** - Get additional seller inputs (legacy)

### 11. **Messaging API** (`messaging`)
- 🔲 **GetMessagingActionsForOrder** - Get messaging actions for order
- 🔲 **ConfirmCustomizationDetails** - Confirm customization details
- 🔲 **CreateConfirmDeliveryDetails** - Create delivery confirmation
- 🔲 **CreateLegalDisclosure** - Create legal disclosure
- 🔲 **CreateNegativeFeedbackRemoval** - Create negative feedback removal
- 🔲 **CreateConfirmServiceDetails** - Create service confirmation
- 🔲 **CreateAmazonMotors** - Create Amazon Motors message
- 🔲 **CreateWarranty** - Create warranty message
- 🔲 **GetAttributes** - Get messaging attributes
- 🔲 **CreateDigitalAccessKey** - Create digital access key
- 🔲 **CreateUnexpectedProblem** - Report unexpected problem

### 12. **Notifications API** (`notifications`)
- 🔲 **GetSubscription** - Get subscription details
- 🔲 **CreateSubscription** - Create subscription
- 🔲 **GetSubscriptionById** - Get subscription by ID
- 🔲 **DeleteSubscriptionById** - Delete subscription
- 🔲 **GetDestinations** - Get notification destinations
- 🔲 **CreateDestination** - Create notification destination
- 🔲 **GetDestination** - Get destination details
- 🔲 **DeleteDestination** - Delete destination

### 13. **Orders API V0** (`ordersV0`) 
- ✅ **GetOrders** - List orders (IMPLEMENTED)
- ✅ **GetOrder** - Get order details (IMPLEMENTED)
- ✅ **GetOrderAddress** - Get order shipping address (IMPLEMENTED)
- ✅ **GetOrderBuyerInfo** - Get order buyer information (IMPLEMENTED)
- ✅ **GetOrderItems** - Get order items (IMPLEMENTED)
- ✅ **GetOrderItemsBuyerInfo** - Get order items buyer info (IMPLEMENTED)

### 14. **Product Pricing API** (`productPricing`)
- 🔲 **GetPricing** - Get pricing for products
- 🔲 **GetCompetitivePricing** - Get competitive pricing
- 🔲 **GetListingOffers** - Get listing offers
- 🔲 **GetItemOffers** - Get item offers
- 🔲 **GetItemOffersBatch** - Get item offers in batch

### 15. **Reports API** (`reports`)
- ✅ **GetReports** - Get report processing status (IMPLEMENTED)
- ✅ **CreateReport** - Create a report (IMPLEMENTED)
- ✅ **GetReport** - Get report details (IMPLEMENTED)
- 🔲 **CancelReport** - Cancel report
- 🔲 **GetReportSchedules** - Get report schedules
- 🔲 **CreateReportSchedule** - Create report schedule
- 🔲 **GetReportSchedule** - Get report schedule details
- 🔲 **CancelReportSchedule** - Cancel report schedule
- ✅ **GetReportDocument** - Get report document (IMPLEMENTED)

### 16. **Sales API** (`sales`)
- ✅ **GetOrderMetrics** - Get order metrics (IMPLEMENTED)

### 17. **Sellers API** (`sellers`)
- 🔲 **GetMarketplaceParticipations** - Get marketplace participations

### 18. **Service API** (`service`)
- 🔲 **GetServiceJobs** - Get service jobs
- 🔲 **GetServiceJobByServiceJobId** - Get service job details
- 🔲 **CancelServiceJobByServiceJobId** - Cancel service job
- 🔲 **CompleteServiceJobByServiceJobId** - Complete service job
- 🔲 **GetServiceJobByServiceJobId** - Get service job by ID
- 🔲 **RescheduleAppointmentForServiceJobByServiceJobId** - Reschedule appointment

### 19. **Shipping API** (`shipping`)
- 🔲 **CreateShipment** - Create shipment
- 🔲 **GetShipment** - Get shipment details
- 🔲 **CancelShipment** - Cancel shipment
- 🔲 **PurchaseLabels** - Purchase shipping labels
- 🔲 **RetrieveShippingLabel** - Retrieve shipping label
- 🔲 **PurchaseShipment** - Purchase shipment
- 🔲 **GetRates** - Get shipping rates
- 🔲 **GetAccount** - Get account information
- 🔲 **GetTrackingInformation** - Get tracking information

### 20. **Small and Light API** (`smallAndLight`)
- 🔲 **GetSmallAndLightEnrollmentBySellerSKU** - Get S&L enrollment by SKU
- 🔲 **PutSmallAndLightEnrollmentBySellerSKU** - Enroll SKU in S&L
- 🔲 **DeleteSmallAndLightEnrollmentBySellerSKU** - Remove SKU from S&L
- 🔲 **GetSmallAndLightEligibilityBySellerSKU** - Check S&L eligibility
- 🔲 **GetSmallAndLightFeePreview** - Get S&L fee preview

### 21. **Solicitations API** (`solicitations`)
- 🔲 **GetSolicitationActionsForOrder** - Get solicitation actions
- 🔲 **CreateProductReviewAndSellerFeedbackSolicitation** - Create review solicitation

### 22. **Uploads API** (`uploads`)
- 🔲 **CreateUploadDestinationForResource** - Create upload destination

---

## Implementation Priority

### **Phase 1: Core Operations** (High Priority)
1. **Orders API** - ✅ COMPLETE
   - ✅ GetOrderAddress 
   - ✅ GetOrderBuyerInfo
   - ✅ GetOrderItems
   - ✅ GetOrderItemsBuyerInfo

2. **Reports API** - ✅ COMPLETE (Essential for business intelligence)
   - ✅ GetReports
   - ✅ CreateReport
   - ✅ GetReport
   - ✅ GetReportDocument

3. **FBA Inventory API** - Critical inventory management
   - 🔲 GetInventorySummaries

4. **Product Pricing API** - Pricing strategy tools
   - 🔲 GetPricing
   - 🔲 GetCompetitivePricing

5. **Finances API** - Financial tracking
   - 🔲 ListFinancialEvents
   - 🔲 ListFinancialEventGroups

### **Phase 2: Enhanced Functionality** (Medium Priority)
6. **Catalog API** - Product discovery
   - 🔲 SearchCatalogItems
   - 🔲 GetCatalogItem

7. **Listings Items API** - Product management
   - 🔲 GetListingsItem
   - 🔲 PutListingsItem

8. **Fees API** - Cost analysis
   - 🔲 GetMyFeesEstimateForSKU
   - 🔲 GetMyFeesEstimateForASIN

9. **Sellers API** - Account information
   - 🔲 GetMarketplaceParticipations

### **Phase 3: Advanced Features** (Lower Priority)
10. **FBA Inbound/Outbound APIs** - Advanced fulfillment
11. **Messaging API** - Customer communication
12. **Notifications API** - Event subscriptions
13. **Shipping/Merchant Fulfillment** - Advanced shipping

---

## Development Guidelines

- Each function should be implemented in its own PR
- Follow existing patterns in `internal/tools/orders.go` and `internal/tools/sales.go`
- Include comprehensive tests for each implementation
- Update `internal/tools/registry.go` to register new tools
- Ensure proper error handling and validation
- Use the basic Client pattern (not ClientWithResponses) to avoid JSON parsing issues