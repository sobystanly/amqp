# Order-Management

The Order Management Service is responsible for handling orders placed by customers,
managing products. This service provides a RESTful API for interacting with orders, customers 
and products.

1. [Usage](#usage)
2. [API Endpoints](#api-endpoints)
3. [Message-Broker](#message-broker)

## Usage
- Ensure that you have followed the setup instructions in the Setup.md in the root folder to set up the 
  Management Service and its dependencies.
- Once the Order Management Service is running, you can interact with it using its RESTful API. The service provides
  endpoints for managing orders, customers, and products

## API ENDPOINTS
The Order Management Service exposes the following endpoints

- **if had more time would've done more proper status code return and validation based on entity. For now for success 
  return 200, 201 and 500 for errors.** 

**Order APIs**

**Create a new order**
```http request
POST http://localhost:8001/orderManagement/order 
```
Payload
```json
  {
  "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce",
  "products": [
    {
      "productId": "e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15",
      "Price": 45,
      "name": "Thermos Flask",
      "productType": "Home & Kitchen",
      "orderedQuantity": 2
    }
  ]
}
```
Response
- Status Code:  201 created
```json
{
    "orderId": "72dd4c34-fa17-11ee-99ad-f40f24119ce9",
    "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce",
    "orderStatus": "placed",
    "paymentStatus": "pending",
    "orderDate": "2024-04-14T00:28:31.336146-04:00",
    "products": [
        {
            "productId": "e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15",
            "name": "Thermos Flask",
            "price": 45,
            "description": "",
            "quantityAvailable": 24,
            "productType": "Home & Kitchen",
            "OrderedQuantity": 2
        }
    ]
}
```
**Get order**

**This API can be extended with many filters and with pagination but for the interest of time 
and simplicity only implementing the basic API that fetch all the orders all the time.**

**If had more time, would've added more filters as follows**
 - Filter by Order Attributes - `order_status`, `customer_id`, `order_date`
 - Filter by order Attributes - `product_id`, `name`, `price`, `product_type`
 - Pagination - `offset`, `limit`, `total records`
 - Sorting - `order_date`, `price`, `product_name`
 - Also may be combining filters with OR and AND, NOT operators etc.

```http request
GET http://localhost:8001localhost:8001/orderManagement/order
```

Response
- Status Code - 200 Ok
```json
[
    {
        "orderId": "72dd4c34-fa17-11ee-99ad-f40f24119ce9",
        "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce",
        "orderStatus": "shipped",
        "paymentStatus": "success",
        "orderDate": "0001-01-01T00:00:00Z",
        "products": [
            {
                "productId": "e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15",
                "name": "Thermos Flask",
                "price": 45,
                "description": "2L keep cold or warm for 48hrs..",
                "quantityAvailable": 23,
                "productType": "Home & Kitchen",
                "OrderedQuantity": 2
            }
        ]
    }
]
```

**Product APIs **

**Add a new Product**

```http request
POST http://localhost:8001/order-management/product
```
Request
```json
  {
  "name": "Boss Speakers",
  "price": 50,
  "description": "Bluetooth Speakers",
  "quantityAvailable": 40,
  "productType": "Electronics"
}
```
Response

```json
{
    "productId": "252e7252-acc3-42fb-b6a7-94ca186b95e8",
    "name": "Boss Speakers",
    "price": 50,
    "description": "Bluetooth Speakers",
    "quantityAvailable": 40,
    "productType": "Electronics",
    "OrderedQuantity": 0
}
```
## Customer API
- For simplicity and interest of time only API to Add a new Customer is added.

Add a new Customer
```http request
POST http://localhost:8001/order-management/customer
```
Request
```json
{
    "firstName": "Emmanuel",
    "lastName": "Joseph",
    "email": "ej@test.com"
}
```

Response
```json
{
    "customerId": "a9d3a80d-6822-44f7-a8d9-61f6f86c779f",
    "firstName": "Emmanuel",
    "lastName": "Joseph",
    "email": "ej@test.com"
}
```
## message-broker
The Order Management service utilizes RabbitMQ message broker, to facilitate communication with the Payment 
processing service. 

This integration streamlines the handling of payment events and enables real-time updates on payment statuses.

It uses the custom AMQP library in this repo to publish payment events to designated queues.
