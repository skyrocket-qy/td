using UnityEngine;

[CreateAssetMenu(fileName = "BleedingEffect", menuName = "Game/Status Effects/Bleeding")]
public class BleedingStatusEffectData : StatusEffectData
{
    public int damagePerTick = 1;

    public override void ApplyEffect(Monster target)
    {
        base.ApplyEffect(target);
        Debug.Log($"Applying Bleeding to {target.name}! {damagePerTick} damage every {tickInterval}s.");
        // Visuals for bleeding could be applied here
    }

    public override void TickEffect(Monster target)
    {
        target.TakeDamage(damagePerTick);
        Debug.Log($"{target.name} takes {damagePerTick} bleeding damage.");
    }

    public override void RemoveEffect(Monster target)
    {
        base.RemoveEffect(target);
        Debug.Log($"{target.name} is no longer bleeding.");
        // Remove bleeding visuals here
    }
}
