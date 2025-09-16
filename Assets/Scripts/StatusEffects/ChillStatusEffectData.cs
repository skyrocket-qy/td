using UnityEngine;

[CreateAssetMenu(fileName = "ChillEffect", menuName = "Game/Status Effects/Chill")]
public class ChillStatusEffectData : StatusEffectData
{
    public float slowAmount = 0.5f; // e.g., 0.5 means 50% slower

    public override void ApplyEffect(Monster target)
    {
        base.ApplyEffect(target);
        Debug.Log($"Chilling {target.name} by {slowAmount * 100}%!");
        // Apply slow to monster's movement speed
        // target.GetComponent<MonsterMovement>().ApplySlow(slowAmount);
    }

    public override void RemoveEffect(Monster target)
    {
        base.RemoveEffect(target);
        Debug.Log($"{target.name} is no longer chilled.");
        // Remove slow from monster's movement speed
        // target.GetComponent<MonsterMovement>().RemoveSlow(slowAmount);
    }
}
