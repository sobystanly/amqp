# Payment-Processing

The payment processing service is responsible for handling payment events emitted by the Order Management service and 
asynchronously communicating back the payment status.

This service integrates with RabbitMQ to listen for payment events and responds with the corresponding payment status 
updates.

## Usage
- Listening for Payment Events
   The order management service emit payment event to payment processing exchange on processPayment queue and this message
   essential info carry out a payment. 
- As per the instruction it simulates the payment as, if the total amount is over $1000 payment is considered failed and
  if the total amount is less than &1000 the payment is considered to be okay