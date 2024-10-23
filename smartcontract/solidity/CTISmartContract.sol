pragma solidity ^0.5.11;
pragma experimental ABIEncoderV2;

contract CTISmartContract {
    // Mapping to store CTI information with CTI ID as the key
    mapping(string => ThreatIntel) public ctis;

    // Array to keep track of the keys in the mapping
    string[] internal keys;

    // Struct for storing Threat Intelligence information
    struct ThreatIntel {
        string id;
        string name;
        string publisher;
        uint256 _type; // Renamed from 'type' to '_type'
        string data;
        string hash;
        uint256 dataSize;
        uint256 value;
        string chainId;
    }

    // Constructor
    constructor() public {
        // Initialization can be done here if needed
    }

    /**
     * Registers a new threat intelligence entry.
     *
     * @param _ctiId CTI's unique identifier.
     * @param _ctiName CTI's name.
     * @param _publisher Publisher of the CTI.
     * @param _type Type of the CTI.
     * @param _data Data associated with the CTI.
     * @param _hash Hash or chain ID of the CTI data.
     * @param _dataSize Size of the CTI data.
     * @param _value Value assessed by the platform.
     * @param _chainId Chain ID where the data is stored.
     * @return success True if the registration was successful, false otherwise.
     */
    function registerCTI(string memory _ctiId, string memory _ctiName, string memory _publisher, uint256 _type,
                         string memory _data, string memory _hash, uint256 _dataSize, uint256 _value,
                         string memory _chainId) public returns (bool success) {
        require(bytes(_ctiId).length > 0 && bytes(ctis[_ctiId].id).length == 0, "CTI ID already exists or is empty");
        ctis[_ctiId] = ThreatIntel(_ctiId, _ctiName, _publisher, _type, _data, _hash, _dataSize, _value, _chainId);
        keys.push(_ctiId);
        return true;
    }

    /**
     * Queries a threat intelligence entry by its ID.
     *
     * @param _ctiId The unique identifier of the CTI.
     * @return The ThreatIntel object if found, null otherwise.
     */
    function queryCTI(string memory _ctiId) public view returns (string memory, string memory, string memory, uint256, string memory, string memory, uint256, uint256, string memory) {
        ThreatIntel storage cti = ctis[_ctiId];
        return (cti.id, cti.name, cti.publisher, cti._type, cti.data, cti.hash, cti.dataSize, cti.value, cti.chainId);
    }

    /**
     * Queries all registered threat intelligence entries.
     *
     * @return A list containing all ThreatIntel objects.
     */
    function queryAllCTIs() public view returns (string[] memory, string[] memory, uint256[] memory, string[] memory, string[] memory, uint256[] memory, uint256[] memory, string[] memory) {
        string[] memory ids = new string[](keys.length);
        string[] memory names = new string[](keys.length);
        uint256[] memory _types = new uint256[](keys.length);
        string[] memory datas = new string[](keys.length);
        string[] memory hashes = new string[](keys.length);
        uint256[] memory dataSizes = new uint256[](keys.length);
        uint256[] memory values = new uint256[](keys.length);
        string[] memory chainIds = new string[](keys.length);

        for (uint256 i = 0; i < keys.length; i++) {
            string memory id = keys[i];
            ThreatIntel storage cti = ctis[id];
            ids[i] = cti.id;
            names[i] = cti.name;
            _types[i] = cti._type;
            datas[i] = cti.data;
            hashes[i] = cti.hash;
            dataSizes[i] = cti.dataSize;
            values[i] = cti.value;
            chainIds[i] = cti.chainId;
        }
        return (ids, names, _types, datas, hashes, dataSizes, values, chainIds);
    }
}