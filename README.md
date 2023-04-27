# owt_contacts_api
## “Contacts API“. It’s a simple API, where a user can get a quick overview over all contacts resources like person and skills. It is implemented by go and uses JWT token for authentication.

The following use cases are implemented:
## UC1
CRUD endpoint for managing contacts. A contact has the following attributes and appropriate validation:

• Firstname

• Lastname

• Fullname

• Address

• Email

• Mobile phone number

## UC2
CRUD endpoint for skills. A contact has multiple skills and a skill belongs to multiple contacts. A skill has the following attributes and appropriate validation:

• Name

• Level (basic, intermediate, advanced, expertise)

## UC3
Document API with Swagger.

## UC4 
The following security aspects are implemented

• Authentication with token

• Authorization: Users can only change their own contact and skills