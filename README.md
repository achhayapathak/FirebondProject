# Firebond Project - Backend Developer Assignment

This repository contains the code for the Backend Developer Assignment of the Firebond Project. The assignment involves building a RESTful API that retrieves exchange rates for cryptocurrencies and fiat currencies, stores the data in a database, integrates with Ethereum blockchain to retrieve account balances, and includes error handling, validation, and testing.

## Table of Contents

- [Minimum Requirements](#minimum-requirements)
- [Installation and Setup](#installation-and-setup)
- [Deployment](#deployment)
- [API Endpoints](#api-endpoints)

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

3. Database Setup:

    1. Sign in to your MongoDB Atlas account or create a new account if you don't have one.

    2. Create a new project and configure the security settings and network access as per your requirements.

    3. Create a database named "Currency_Exchange".

    4. Get the MongoDB connection URI provided by MongoDB Atlas for this database.

    5. Set up the connection URI in the `.env` file. Follow the instructions in the [.env File Setup](#env-file-setup) section.

4. Setup the environment variables:

    1. Create a new file named `.env` in the root directory of the project.

    2. Open the `.env` file and add the following lines:

        ```bash
        API_KEY=your_cryptocompare_api_key
        MONGO_URI=your_mongo_connection_uri
        INFURA_URI=your_infura_uri
        ```

        Replace <your_cryptocompare_api_key> with your actual CryptoCompare API key for fetching exchange rates.
        Replace <your_mongo_connection_uri> with the MongoDB connection URI obtained in the previous step.
        Replace <your_infura_uri> with your actual Infura URI for fetching balance from ethereum address.

5. Start the server:

    ```bash
    go run main.go
    ```

6. The API should now be running locally at `http://localhost:8080`.


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
