using UnityEngine;

[CreateAssetMenu(fileName = "PoisonEffect", menuName = "Game/Status Effects/Poison")]
public class PoisonStatusEffectData : StatusEffectData
{
    public int damagePerTick = 1;

    public override void ApplyEffect(Monster target)
    {
        base.ApplyEffect(target);
        Debug.Log($"{target.name} is now poisoned! {damagePerTick} damage every {tickInterval}s.");
        // Visuals for poison could be applied here
    }

    public override void TickEffect(Monster target)
    {
        target.TakeDamage(damagePerTick);
        Debug.Log($"{target.name} takes {damagePerTick} poison damage.");
    }

    public override void RemoveEffect(Monster target)
    {
        base.RemoveEffect(target);
        Debug.Log($"{target.name} is no longer poisoned.");
        // Remove poison visuals here
    }
}
