using UnityEngine;
using UnityEngine.AI; // Required for NavMeshAgent

/*
 * To make this MonsterMovement script work with NavMesh in your Unity project, follow these steps:
 *
 * 1. Enable AI Navigation Package:
 *    - In Unity, go to `Window` > `Package Manager`.
 *    - Make sure "Unity Registry" is selected in the dropdown.
 *    - Search for "AI Navigation" and install it.
 *
 * 2. Mark Walkable Surfaces:
 *    - Select all the GameObject(s) in your scene that represent the ground or surfaces where monsters should be able to walk.
 *    - In the Inspector window, find the "Navigation Static" checkbox (you might need to enable "Static" first).
 *    - Check the "Navigation Static" checkbox.
 *    - In the dropdown that appears, select "Walkable".
 *
 * 3. Mark Obstacles (Barriers):
 *    - Select all the GameObject(s) in your scene that represent barriers or obstacles that monsters should avoid.
 *    - In the Inspector window, check the "Navigation Static" checkbox.
 *    - In the dropdown that appears, select "Not Walkable" or "Obstacle".
 *
 * 4. Bake the NavMesh:
 *    - Go to `Window` > `AI` > `Navigation`. This will open the Navigation window.
 *    - Go to the "Bake" tab.
 *    - Adjust the "Agent Radius", "Agent Height", "Max Slope", and "Drop Height" settings as needed for your monster's size and movement capabilities.
 *    - Click the "Bake" button. Unity will generate a blue overlay on your walkable surfaces, representing the NavMesh.
 *
 * 5. Add `NavMeshAgent` to Monster:
 *    - Select your monster GameObject (the one with the `Monster` and `MonsterMovement` scripts).
 *    - In the Inspector, click "Add Component" and search for "Nav Mesh Agent". Add this component.
 *    - Adjust the "Speed", "Angular Speed", "Acceleration", and "Stopping Distance" properties of the `NavMeshAgent`
 *      to control your monster's movement behavior. The `speed` variable in `MonsterMovement.cs` is no longer used;
 *      the `NavMeshAgent`'s speed will control the movement.
 */
[RequireComponent(typeof(NavMeshAgent))] // Ensures a NavMeshAgent is present
public class MonsterMovement : MonoBehaviour
{
    public Vector3 start;
    public Vector3 end;
    // Speed will now be controlled by NavMeshAgent.speed

    private NavMeshAgent agent;

    void Awake()
    {
        agent = GetComponent<NavMeshAgent>();
        if (agent == null)
        {
            Debug.LogError("NavMeshAgent component not found on this GameObject!");
        }
    }

    void Start()
    {
        // Set the agent's starting position
        if (agent != null)
        {
            agent.Warp(start); // Teleport agent to start position
            agent.SetDestination(end); // Set the target destination
        }
    }

    void Update()
    {
        // Check if the agent has reached its destination
        if (agent != null && !agent.pathPending && agent.remainingDistance <= agent.stoppingDistance)
        {
            if (!agent.hasPath || agent.velocity.sqrMagnitude == 0f)
            {
                Debug.Log("Monster reached the end! Reducing player lives.");
                GameManager gameManager = FindAnyObjectByType<GameManager>();
                if (gameManager != null)
                {
                    gameManager.playerLives--;
                    Debug.Log("Player Lives: " + gameManager.playerLives);
                }
                Destroy(gameObject);
            }
        }
    }
}
