using UnityEngine;
using System.Collections.Generic;

[CreateAssetMenu(fileName = "NewPassiveTree", menuName = "Game/Progression/Passive Tree")]
public class PassiveTreeData : ScriptableObject
{
    public string treeName = "New Passive Tree";
    public List<PassiveAbilityData> passiveAbilities = new List<PassiveAbilityData>(); // All abilities in this tree

    [Header("Runtime Data (for a specific hero)")]
    // IMPORTANT: These fields (totalPassivePoints, unlockedAbilities) are meant to be
    // managed PER HERO INSTANCE. ScriptableObjects are shared assets.
    // To make this work correctly for multiple heroes, you should:
    // 1. Create a runtime instance of PassiveTreeData for each hero (e.g., by instantiating it).
    // 2. Or, store these runtime values (points, unlocked abilities) directly on the Hero component
    //    and pass them to methods like TryUnlockAbility.
    // For simplicity in this example, we're showing them here, but be aware of the ScriptableObject's shared nature.
    public int totalPassivePoints = 0; // Points available to spend on this tree
    public List<PassiveAbilityData> unlockedAbilities = new List<PassiveAbilityData>(); // Abilities currently unlocked

    // In a more complex system, this might be a graph structure
    // For now, it's a simple list of abilities that can be applied.

    // Placeholder method to try and unlock an ability
    public bool TryUnlockAbility(PassiveAbilityData abilityToUnlock)
    {
        if (abilityToUnlock == null) return false;

        // Check if already unlocked
        if (unlockedAbilities.Contains(abilityToUnlock))
        {
            Debug.LogWarning($"{abilityToUnlock.abilityName} is already unlocked.");
            return false;
        }

        // Check if enough passive points are available
        if (totalPassivePoints < abilityToUnlock.unlockCost)
        {
            Debug.LogWarning($"Not enough passive points to unlock {abilityToUnlock.abilityName}. Required: {abilityToUnlock.unlockCost}, Available: {totalPassivePoints}");
            return false;
        }

        // Check prerequisites
        foreach (PassiveAbilityData prerequisite in abilityToUnlock.prerequisites)
        {
            if (!unlockedAbilities.Contains(prerequisite))
            {
                Debug.LogWarning($"Prerequisite {prerequisite.abilityName} not met for {abilityToUnlock.abilityName}.");
                return false;
            }
        }

        // If all checks pass, unlock the ability
        unlockedAbilities.Add(abilityToUnlock);
        totalPassivePoints -= abilityToUnlock.unlockCost;
        Debug.Log($"Successfully unlocked {abilityToUnlock.abilityName}. Remaining points: {totalPassivePoints}");
        return true;
    }

    // Method to reset the tree (for a specific hero)
    public void ResetTree(int startingPoints)
    {
        unlockedAbilities.Clear();
        totalPassivePoints = startingPoints;
        Debug.Log($"Passive tree reset. Starting points: {startingPoints}");
    }
}
