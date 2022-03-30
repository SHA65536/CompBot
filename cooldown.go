package main

import (
	"fmt"
	"time"
)

// CooldownManager controls the different cooldowns
type CooldownManager struct {
	Cooldowns map[string]*Cooldown
}

func MakeCoolownManager() *CooldownManager {
	return &CooldownManager{
		Cooldowns: make(map[string]*Cooldown),
	}
}

// NewCooldown creates a new cooldown, and sets it's threshold,
// If the cooldown already exists, returns an error.
func (cm *CooldownManager) NewCooldown(name string, threshold time.Duration) error {
	if _, ok := cm.Cooldowns[name]; ok {
		return fmt.Errorf("cooldown: cooldown %s already exists", name)
	}
	cm.Cooldowns[name] = &Cooldown{
		Name:      name,
		Times:     make(map[string]time.Time),
		Threshold: threshold,
	}
	return nil
}

// IsAllowed checks if an object is still in cooldown,
// Errors if cooldown does not exist.
// Returns True if object not in cooldown, or already past it.
func (cm *CooldownManager) IsAllowed(name, id string) (bool, error) {
	if cd, ok := cm.Cooldowns[name]; ok {
		return cd.IsAllowed(id), nil
	}
	return false, fmt.Errorf("cooldown: cooldown %s does not exist", name)
}

// SetObject sets an object's timestamp to time.Now().UTC()
// Errors if cooldown does not exist.
func (cm *CooldownManager) SetObject(name, id string) error {
	if cd, ok := cm.Cooldowns[name]; ok {
		cd.SetCooldown(id)
		return nil
	}
	return fmt.Errorf("cooldown: cooldown %s does not exist", name)
}

// SetThreshold sets a cooldown's threshold to the duration given
// Errors if cooldown does not exist.
func (cm *CooldownManager) SetThreshold(name string, threshold time.Duration) error {
	if cd, ok := cm.Cooldowns[name]; ok {
		cd.Threshold = threshold
		return nil
	}
	return fmt.Errorf("cooldown: cooldown %s does not exist", name)
}

// Cleans all of the cooldowns from expired objects
func (cm *CooldownManager) CleanCooldowns() {
	for _, cd := range cm.Cooldowns {
		cd.Clean()
	}
}

type Cooldown struct {
	Name      string               // Cooldown Name
	Times     map[string]time.Time // Last command map
	Threshold time.Duration        // Duration of allowed command
}

func (c *Cooldown) IsAllowed(user string) bool {
	if timestamp, ok := c.Times[user]; ok {
		if time.Now().UTC().Before(timestamp.Add(c.Threshold)) {
			return false
		}
	}
	return true
}

func (c *Cooldown) SetCooldown(user string) {
	c.Times[user] = time.Now().UTC()
}

func (c *Cooldown) Clean() {
	for id, timestamp := range c.Times {
		if time.Now().UTC().After(timestamp.Add(c.Threshold)) {
			delete(c.Times, id)
		}
	}
}
