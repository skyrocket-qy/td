using UnityEngine;
using System.Collections.Generic;

public class StatusEffectManager : MonoBehaviour
{
    // Nested class to hold runtime data for an active status effect
    private class ActiveStatusEffect
    {
        public StatusEffectData effectData;
        public float remainingDuration;
        public float lastTickTime;

        public ActiveStatusEffect(StatusEffectData data)
        {
            effectData = data;
            remainingDuration = data.duration;
            lastTickTime = Time.time;
        }
    }

    private List<ActiveStatusEffect> activeEffects = new List<ActiveStatusEffect>();
    private Monster _monster; // Reference to the monster this manager is on

    void Awake()
    {
        _monster = GetComponent<Monster>();
        if (_monster == null)
        {
            Debug.LogError("StatusEffectManager requires a Monster component on the same GameObject.");
            enabled = false; // Disable if no Monster component
        }
    }

    void Update()
    {
        if (_monster == null) return;

        List<ActiveStatusEffect> effectsToRemove = new List<ActiveStatusEffect>();

        foreach (ActiveStatusEffect activeEffect in activeEffects)
        {
            activeEffect.remainingDuration -= Time.deltaTime;

            // Handle ticking effects (e.g., DoT)
            if (activeEffect.effectData.tickInterval > 0 && Time.time >= activeEffect.lastTickTime + activeEffect.effectData.tickInterval)
            {
                activeEffect.effectData.TickEffect(_monster);
                activeEffect.lastTickTime = Time.time;
            }

            if (activeEffect.remainingDuration <= 0)
            {
                effectsToRemove.Add(activeEffect);
            }
        }

        foreach (ActiveStatusEffect effectToRemove in effectsToRemove)
        {
            effectToRemove.effectData.RemoveEffect(_monster);
            activeEffects.Remove(effectToRemove);
        }
    }

    public void ApplyStatusEffect(StatusEffectData effectData)
    {
        if (_monster == null || effectData == null) return;

        // Check if effect is already active and refresh/stack if needed
        ActiveStatusEffect existingEffect = activeEffects.Find(e => e.effectData == effectData);
        if (existingEffect != null)
        {
            // For simplicity, just refresh duration. Could add stacking logic here.
            existingEffect.remainingDuration = effectData.duration;
            Debug.Log($"Refreshed {effectData.effectName} on {_monster.name}.");
        }
        else
        {
            ActiveStatusEffect newActiveEffect = new ActiveStatusEffect(effectData);
            activeEffects.Add(newActiveEffect);
            effectData.ApplyEffect(_monster);
        }
    }

    // Helper to check if a specific effect is active
    public bool HasStatusEffect(StatusEffectData effectData)
    {
        return activeEffects.Exists(e => e.effectData == effectData);
    }

    // Helper to get an active effect (e.g., to modify its properties)
    public StatusEffectData GetActiveStatusEffect(StatusEffectData effectData)
    {
        return activeEffects.Find(e => e.effectData == effectData)?.effectData;
    }
}
