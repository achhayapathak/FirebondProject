# Firebond Project - Backend Developer Assignment

This repository contains the code for the Backend Developer Assignment of the Firebond Project. The assignment involves building a RESTful API that retrieves exchange rates for cryptocurrencies and fiat currencies, stores the data in a database, integrates with Ethereum blockchain to retrieve account balances, and includes error handling, validation, and testing.

## Table of Contents

- [Minimum Requirements](#minimum-requirements)
- [Installation and Setup](#installation-and-setup)
- [API Endpoints](#api-endpoints)
- [Deployment](#deployment)

## Minimum Requirements

The assignment requires implementing the following features:

1. **Data Acquisition:** Fetch exchange rates for cryptocurrencies and fiat currencies from an existing API (e.g., CoinGecko or CryptoCompare) periodically.

2. **Database:** Implement a database to store the exchange rate data efficiently. It can be an SQL or NoSQL database based on your preference.

3. **API Endpoints:** Develop the following RESTful API endpoints:

   - `GET /rates/{cryptocurrency}/{fiat}`: Returns the current exchange rate between the specified cryptocurrency and fiat currency.
   - `GET /rates/{cryptocurrency}`: Returns the current exchange rates between the specified cryptocurrency and all supported fiat currencies.
   - `GET /rates`: Returns the current exchange rates for all supported cryptocurrency-fiat pairs.
   - `GET /rates/history/{cryptocurrency}/{fiat}`: Returns the exchange rate history between the specified cryptocurrency and fiat currency for the past 24 hours.

4. **Web3 Integration:** Use the web3 library to retrieve the current balance of a specific Ethereum address. Expose this functionality through the API endpoint `GET /balance/{address}`.

5. **Error Handling and Validation:** Implement proper error handling and validation for request data. Return appropriate error messages for invalid requests and handle exceptions gracefully.

6. **Testing:** Write unit and integration tests to verify the correctness of the API.

## Installation and Setup

To run the API locally, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/achhayapathak/FirebondProject.git

2. Install the required dependencies:

    ```bash
    go mod download
    ```

3. Start the server:

    ```bash
    go run main.go
    ```

4. The API should now be running locally at `http://localhost:8080`.


## Deployment

The backend has been deployed on [https://firebondproject.onrender.com/](https://firebondproject.onrender.com/).

You can access the deployed backend at [https://firebondproject.onrender.com/](https://firebondproject.onrender.com/).


## API Endpoints

The following API endpoints are available:

- `GET /rates/{cryptocurrency}/{fiat}`
- `GET /rates/{cryptocurrency}`
- `GET /rates`
- `GET /rates/history/{cryptocurrency}/{fiat}`
- `GET /balance/{address}`

For any further questions or issues, please contact me at [achhayapathak11@gmail.com](mailto:achhayapathak11@gmail.com).
