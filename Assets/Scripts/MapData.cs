using UnityEngine;
using System.Collections.Generic;

[CreateAssetMenu(fileName = "NewMapData", menuName = "Game/Map Data")]
public class MapData : ScriptableObject
{
    public string mapName = "New Map";
    public Vector3 monsterStartPoint = new Vector3(-5f, 0f, 0f);
    public Vector3 monsterEndPoint = new Vector3(5f, 0f, 0f);
    public List<Vector3> barrierPositions = new List<Vector3>();
    public List<Vector3> heroPlacementSpots = new List<Vector3>();
    public int heroLimit = 1; // Max number of heroes allowed on this map
}
