package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const userIdContextKey = contextKey("userId")
