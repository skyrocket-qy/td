using UnityEngine;
using System.Collections.Generic; // Added for List

[CreateAssetMenu(fileName = "NewPassiveAbility", menuName = "Game/Progression/Passive Ability")]
public class PassiveAbilityData : ScriptableObject
{
    public string abilityName = "New Passive Ability";
    [TextArea] public string description = "Description of the passive ability.";
    public Sprite icon; // Optional: for UI display

    [Header("Tree Properties")]
    // Drag other PassiveAbilityData assets here that must be unlocked before this one.
    public List<PassiveAbilityData> prerequisites = new List<PassiveAbilityData>(); 
    public int unlockCost = 1; // The number of passive points required to unlock this ability.
    public Vector2 positionInTree; // Used for visual layout in a potential editor tool or UI.

    [Header("Ability Effects")]
    // These fields define the effects this passive ability grants.
    // The Hero script will read these values and apply them to its stats.
    public int bonusDamage = 0;
    public float bonusAttackRange = 0f;
    public float bonusHealthMultiplier = 1f; // e.g., 1.1 for +10% health

    // You can add more specific effects here, or use a more complex system
    // For now, this is just data. The Hero script would interpret and apply these.
}
