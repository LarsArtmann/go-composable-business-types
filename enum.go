package cbt

//go:generate go tool go-enum --marshal --names --values --mustparse

// ActorKind represents the type of actor performing an action.
// ENUM(User, Bot, System, Service)
type ActorKind uint8

// Locale represents a language-region locale.
// ENUM(en_US, en_GB, de_DE, fr_FR, es_ES, it_IT, ja_JP, zh_CN)
type Locale string

// Priority represents task/issue priority levels.
// ENUM(Low, Medium, High, Critical)
type Priority uint8

// Status represents common entity lifecycle status.
// ENUM(Draft, Active, Paused, Archived, Deleted)
type Status uint8

// Trigger represents what caused a DataPoint to be created.
// ENUM(
//
//	Manual,    // Direct user action
//	Scheduled, // Time-based trigger (cron, delay)
//	Webhook,   // External system via webhook
//	Import,    // Bulk data import
//	Migration, // Data migration
//	System,    // Automatic system action
//	Correction // Correction of previous data
//
// )
type Trigger uint8
