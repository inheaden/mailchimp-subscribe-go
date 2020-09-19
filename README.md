# Mailchimp Subscribe Proxy

This is a simple proxy for the mailchimp subscribe api written in Go.

We use it internally to build 'Subscribe to newsletter' functionalities to public websites.

## Instructions

1. Copy `.env.example` to `.env` or make sure to set `API_KEY`, `LIST_ID`, `API_SERVER` environment variables.
2. Run `go run main.go mc-api.go`

There is also a `Dockerfile` provided.

## API

The service offers two APIs:

### `POST /`

Subscribes the provided email using the [Add or update list member](https://mailchimp.com/developer/api/marketing/list-members/add-or-update-list-member/) API of mailchimp.

#### Body

```json
{
  "email": "EMAIL_TO_SUBSCRIBE"
}
```

### `GET /`

Can be used to check status of the API, returns `{"message": "UP"}`.

## License

This project is licensed under the [MIT](LICENSE) license.
