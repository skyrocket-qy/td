using UnityEngine;

[CreateAssetMenu(fileName = "IgniteEffect", menuName = "Game/Status Effects/Ignite")]
public class IgniteStatusEffectData : StatusEffectData
{
    public int damagePerTick = 2;

    public override void ApplyEffect(Monster target)
    {
        base.ApplyEffect(target);
        Debug.Log($"Igniting {target.name}! {damagePerTick} damage every {tickInterval}s.");
        // Visuals for ignite could be applied here
    }

    public override void TickEffect(Monster target)
    {
        target.TakeDamage(damagePerTick);
        Debug.Log($"{target.name} takes {damagePerTick} ignite damage.");
    }

    public override void RemoveEffect(Monster target)
    {
        base.RemoveEffect(target);
        Debug.Log($"{target.name} is no longer ignited.");
        // Remove ignite visuals here
    }
}
