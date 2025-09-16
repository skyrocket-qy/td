using UnityEngine;

[CreateAssetMenu(fileName = "ShockEffect", menuName = "Game/Status Effects/Shock")]
public class ShockStatusEffectData : StatusEffectData
{
    public float damageTakenIncrease = 0.25f; // e.g., monster takes 25% more damage

    public override void ApplyEffect(Monster target)
    {
        base.ApplyEffect(target);
        Debug.Log($"Shocking {target.name}! Takes {damageTakenIncrease * 100}% more damage.");
        // Modify monster's damage resistance/vulnerability
        // target.damageResistanceMultiplier -= damageTakenIncrease;
    }

    public override void RemoveEffect(Monster target)
    {
        base.RemoveEffect(target);
        Debug.Log($"{target.name} is no longer shocked.");
        // Revert monster's damage resistance/vulnerability
        // target.damageResistanceMultiplier += damageTakenIncrease;
    }
}
