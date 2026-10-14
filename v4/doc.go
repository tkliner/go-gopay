// Package v4 implements the GoPay Payments API v4 client.
//
// It is intentionally separate from the root gopay package (REST API v3 / hosted
// gateway). Existing v3 integrations keep using github.com/tkliner/go-gopay;
// import this package only when you need checkout/charge flows (Apple Pay,
// Google Pay, encrypted card, card tokens).
//
// Documentation: https://api-docs.gopay.com/
package v4
