# Final-Project-JCC-Golang-2022
This is a Ecommerce API made with Golang

## API Policy
- User can register as Guest and then choose to Register as a Seller or Customer
- User can't register as Admin
- Category List are open for everyone, but only admin could add, edit and delete it
- Customer list are open for everyone
- Customer could only edit and delete their own data. When customer delete their data their role will be reverted back to Guest. Same goes for Seller
- Customer can't register again as another role. Same goes for Seller
- Product list are open for everyone, but only Seller could add, edit and delete it
- Seller could only edit and delete their own product
- Admin could also delete a product
- Review List are open for everyone, but only Customer could add new review
- Customer could only edit and delete review made by them
- Transaction list are only open for Admin
- Seller and Customer could see transaction related to them
- Only related Seller and Customer could edit a transaction data
- User can update their own password
- User list are open for everyone (for now)
- User can delete their own data. When they do that their Customer or Seller data will also be deleted

