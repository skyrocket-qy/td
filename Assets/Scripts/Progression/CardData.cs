using UnityEngine;

[CreateAssetMenu(fileName = "NewCard", menuName = "Game/Card Data")]
public class CardData : ScriptableObject
{
    public string cardName = "New Card";
    [TextArea(3, 5)]
    public string description = "A brief description of the card's effect.";
    public Sprite icon; // Visual representation of the card
    public CardEffectType effectType; // Enum to categorize effects (e.g., StatBoost, NewSkill, Passive)
    public float effectValue; // Generic value for effects (e.g., +5 damage, +10% health)

    // You might add more specific fields depending on your effects:
    // public SkillData skillToGrant;
    // public PassiveAbilityData passiveToGrant;
    // public HeroClassData classToUnlock;

    public void ApplyEffect(Player player)
    {
        // This method will contain the logic for what the card does
        // You'll implement this based on effectType and other fields
        Debug.Log($"Applying effect of card: {cardName}");
        switch (effectType)
        {
            case CardEffectType.StatBoost:
                // Example: player.IncreaseDamage(effectValue);
                break;
            case CardEffectType.NewSkill:
                // Example: player.LearnSkill(skillToGrant);
                break;
            // ... other cases
        }
    }
}

public enum CardEffectType
{
    None,
    StatBoost,
    NewSkill,
    PassiveAbility,
    // Add more as needed
}
