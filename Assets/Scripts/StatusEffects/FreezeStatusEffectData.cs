using UnityEngine;

[CreateAssetMenu(fileName = "FreezeEffect", menuName = "Game/Status Effects/Freeze")]
public class FreezeStatusEffectData : StatusEffectData
{
    public override void ApplyEffect(Monster target)
    {
        base.ApplyEffect(target);
        Debug.Log($"Freezing {target.name}!");
        // Disable monster movement and animation
        // target.GetComponent<MonsterMovement>().DisableMovement();
        // target.GetComponent<Animator>().speed = 0;
    }

    public override void RemoveEffect(Monster target)
    {
        base.RemoveEffect(target);
        Debug.Log($"{target.name} is no longer frozen.");
        // Re-enable monster movement and animation
        // target.GetComponent<MonsterMovement>().EnableMovement();
        // target.GetComponent<Animator>().speed = 1;
    }
}
